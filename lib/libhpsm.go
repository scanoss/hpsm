package main

/*
 struct ranges{
	char *local;
	char *remote;
	char *matched;
};

*/
import "C"
import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"unsafe"

	m "scanoss.com/hpsm/API/go"
	"scanoss.com/hpsm/model"
	proc "scanoss.com/hpsm/pkg"
	u "scanoss.com/hpsm/utils"
)

/**
go build -o libhpsm.so  -buildmode=c-shared libhpsm.go
gcc -v client.c -o client ./libhpsm.so
*/
// Go_HandleData converts a unsigned char [] C array to an array
// of GO bytes

func Go_handleData(data *C.uchar, length C.int) []byte {
	return C.GoBytes(unsafe.Pointer(data), C.int(length))
}

// Get the file contents of a given url name and place it on
// file
func GetFileContent(url string, filepath string) error {
	// run shell `wget URL -O filepath`

	cmd := exec.Command("wget", url, "-O", filepath, "-T", "10")
	return cmd.Run()
}

//export HashFileContents
func HashFileContents(data *C.char) *C.char {
	dataArray := C.GoString(data)
	var out string = "hpsm="
	hashLocal := proc.GetLineHashesFromSource(dataArray)

	for i := range hashLocal {
		a := fmt.Sprintf("%02x", hashLocal[i])
		out += a
	}
	return C.CString(out)
}

//export HPSM
func HPSM(data *C.char, md5 *C.char) C.struct_ranges {
	dataArray := C.GoString(data)
	var crcSource []byte

	for i := 0; i < len(dataArray)-2; i += 2 {
		var thisString string
		thisString = dataArray[i : i+2]
		thisByte, err := strconv.ParseInt(thisString, 16, 9)
		if err == nil {
			crcSource = append(crcSource, byte(thisByte))
		}
	}

	MD5 := C.GoString(md5)
	//Remote access
	strLinesGo := ""
	ossLinesGo := ""
	totalLines := len(crcSource)
	snippets := localProcessHPSM(crcSource, MD5, 5)
	//replace the above line if API processing is needed
	//snippets := remoteProcessHPSM(crcSource, MD5, 5)
	matchedLines := 0
	for i := range snippets {
		var ossRange string
		var srcRange string
		matchedLines += (snippets[i].LEnd - snippets[i].LStart)
		srcRange = fmt.Sprintf("%d-%d", snippets[i].LStart+1, snippets[i].LEnd+1)
		ossRange = fmt.Sprintf("%d-%d", snippets[i].RStart+1, snippets[i].REnd+1)
		strLinesGo += srcRange
		ossLinesGo += ossRange
		if i < len(snippets)-1 {
			strLinesGo += ", "
			ossLinesGo += ", "
		}

	}
	mLines := ""
	if totalLines == 0 {
		mLines = fmt.Sprint("0%")
	} else {
		mLines = fmt.Sprintf("%d%%", matchedLines*100/totalLines)
	}
	var lines C.struct_ranges
	lines.local = ((*C.char)(C.CString(strLinesGo)))
	lines.remote = ((*C.char)(C.CString(ossLinesGo)))
	lines.matched = ((*C.char)(C.CString(mLines)))

	return lines

}

//export ProcessHPSM
func ProcessHPSM(data *C.uchar, length C.int, md5 *C.char) C.struct_ranges {
	dataArray := Go_handleData(data, length)
	MD5 := C.GoString(md5)
	//Remote access
	strLinesGo := ""
	ossLinesGo := ""
	snippets := localProcessHPSM(dataArray, MD5, 5)
	for i := range snippets {
		var ossRange string
		var srcRange string
		srcRange = fmt.Sprintf("%d-%d", snippets[i].LStart, snippets[i].LEnd)
		ossRange = fmt.Sprintf("%d-%d", snippets[i].RStart, snippets[i].REnd)
		strLinesGo += ossRange
		ossLinesGo += srcRange
		if i < len(snippets)-1 {
			strLinesGo += ", "
			ossLinesGo += ", "
		}

	}

	var lines C.struct_ranges
	lines.local = ((*C.char)(C.CString(strLinesGo)))
	lines.remote = ((*C.char)(C.CString(ossLinesGo)))

	return lines

}

func remoteProcessHPSM(local []uint8, remoteMd5 string, Threshold uint32) []model.Range {

	var req []m.HpsmReqItem
	var item m.HpsmReqItem
	var outRange []model.Range
	item.MD5 = remoteMd5
	item.Hashes = local
	req = append(req, item)

	// Create the HPSM Req JSON

	out, _ := json.Marshal(req)
	fmt.Println(string(out))
	//Request HPSM via CURL
	hpsm := u.RequestHPSM("http://ns3193417.ip-152-228-225.eu:8081", string(out))
	//return scan results + HPSM
	var resp []m.HpsmRespItem
	_ = json.Unmarshal(hpsm, &resp)
	for i := range resp {
		snippets := resp[i].Snippets
		var r model.Range

		for s := range snippets {
			r.RStart = int(snippets[s].Remote.Start)
			r.REnd = int(snippets[s].Remote.End)
			r.LStart = int(snippets[s].Local.Start)
			r.LEnd = int(snippets[s].Local.End)
			outRange = append(outRange, r)
		}
	}
	return outRange
}

func localProcessHPSM(local []uint8, remoteMd5 string, Threshold uint32) []model.Range {
	//Remote access to API

	MD5 := remoteMd5
	srcEndpoint := os.Getenv("SRC_URL")
	if srcEndpoint == "" {
		srcEndpoint = "https://osskb.org/api/file_contents/"
	}
	err := GetFileContent(srcEndpoint+MD5, "/tmp/"+MD5)
	if err == nil {
		hashRemote := proc.GetLineHashes("/tmp/" + MD5)
		u.Rm("/tmp/" + MD5)
		return proc.Compare(local, hashRemote, uint32(5))
	} else {
		return []model.Range{}
	}

}

func main() {}
