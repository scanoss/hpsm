// SPDX-License-Identifier: GPL-2.0-or-later
/*
 * Copyright (C) 2018-2022 SCANOSS.COM
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 2 of the License, or
 * (at your option) any later version.
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

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
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"unsafe"

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

	args := []string{url, "-O", filepath, "-T", "10"}

	// Set X-Session header if SCANOSS_API_KEY is present
	apiKey := os.Getenv("SCANOSS_API_KEY")
	if apiKey != "" {
		args = append(args, "--header=X-Session: "+apiKey)
	}

	cmd := exec.Command("wget", args...)
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

func localProcessHPSM(local []uint8, remoteMd5 string, Threshold uint32) []model.Range {
	//Remote access to API

	MD5 := remoteMd5
	srcEndpoint := os.Getenv("SCANOSS_FILE_CONTENTS_URL")
	if srcEndpoint == "" {
		srcEndpoint = "localhost/api/file_contents/"
	}
	err := GetFileContent(srcEndpoint+MD5, "/tmp/"+MD5)

	if err == nil {
		hashRemote := proc.GetLineHashes("/tmp/" + MD5)
		if len(hashRemote) <= 5 {
			r := model.Range{LStart: -1, LEnd: -1, RStart: -1, REnd: -1}
			return []model.Range{r}
		}

		u.Rm("/tmp/" + MD5)
		return proc.Compare(local, hashRemote, uint32(5))
	} else {
		r := model.Range{LStart: -1, LEnd: -1, RStart: -1, REnd: -1}
		return []model.Range{r}
	}

}

func main() {}
