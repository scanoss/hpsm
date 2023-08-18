package main

/**

This file wraps HPSM functions to be used from another programming language
creating a shared library (.so or .dll) by downloading the remote file and
doing local Calculation of HPSM

Exported functions:
- HPSM: Formats an input string into []byte and send it with a MD5 of the OSS file.
  Computing HPSM is locally by downloading the remote OSS file and later comparisson
- HashFileContents: This auxiliar function is used from a different language to calculate hashes

To build the library:
go build -o libhpsm.so  -buildmode=c-shared libhpsm.go
In order to buld a C client with the library
gcc -v client.c -o client ./libhpsm.so

By default, OSS files are downlaoded from osskb.org
Server address can be set by using Enviroment variable SRC_URL
*/

/*
 struct ranges{
	char *local;
	char *remote;
	char *matched;
};
*/
import "C"
import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"unsafe"

	"scanoss.com/hpsm/model"
	proc "scanoss.com/hpsm/pkg"
	u "scanoss.com/hpsm/utils"
)

// Go_HandleData converts a unsigned char [] C array to an array
// of GO bytes

func Go_handleData(data *C.uchar, length C.int) []byte {
	return C.GoBytes(unsafe.Pointer(data), C.int(length))
}

// Get the file contents of a given its url name and place it on
// dst location
func GetFileContent(url string, dst string) error {
	// run shell `wget URL -O filepath`

	cmd := exec.Command("wget", url, "-O", dst, "-T", "10")
	return cmd.Run()
}

//Auxiliar function to calculate hashes from a source code
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

//Calls HPSM local caclulation. The fist parameter is a HPSM definition (hpsm=zzyyww...) and
// the second is the key of the file to be compared. Returns a struct interpreted by C containting Ranges and matched percentage
//export HPSM
func HPSM(data *C.char, md5 *C.char) C.struct_ranges {
	dataArray := C.GoString(data)
	var crcSource []byte

	for i := 0; i < len(dataArray)-2; i += 2 {
		thisString := dataArray[i : i+2]
		thisByte, err := strconv.ParseInt(thisString, 16, 9)
		if err == nil {
			crcSource = append(crcSource, byte(thisByte))
		}
	}

	MD5 := C.GoString(md5)

	strLinesGo := ""
	ossLinesGo := ""
	totalLines := len(crcSource)
	snippets := localProcessHPSM(crcSource, MD5, 5)

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
		mLines = "0%"
	} else {
		mLines = fmt.Sprintf("%d%%", matchedLines*100/totalLines)
	}
	var lines C.struct_ranges
	lines.local = ((*C.char)(C.CString(strLinesGo)))
	lines.remote = ((*C.char)(C.CString(ossLinesGo)))
	lines.matched = ((*C.char)(C.CString(mLines)))

	return lines

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
