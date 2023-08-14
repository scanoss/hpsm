package process

import (
	"bytes"
	"os"
	"strings"
	"sync"

	"github.com/sigurn/crc8"
	model "scanoss.com/hpsm/model"
)

// Hashes a file
// Calculates CRC-8 of each line contained on a file
func GetLineHashes(fileName string) []uint8 {
	src, err := os.ReadFile(fileName)
	var srcChk []uint8

	if err == nil {

		table := crc8.MakeTable(crc8.CRC8_MAXIM)

		linesSrc := strings.Split(string(src), "\n")
		for line := range linesSrc {
			if linesSrc[line] == "" {
				srcChk = append(srcChk, 0xFF)
				continue
			}
			checksum := crc8.Checksum([]byte(Normalize(linesSrc[line])), table)
			srcChk = append(srcChk, checksum)
		}
	}
	return srcChk

}

// Hashes a file
// Calculates CRC-8 of each line contained on a source code string
func GetLineHashesFromSource(src string) []uint8 {

	var srcChk []uint8

	table := crc8.MakeTable(crc8.CRC8_MAXIM)
	linesSrc := strings.Split(string(src), "\n")

	for line := range linesSrc {
		if linesSrc[line] == "" {
			srcChk = append(srcChk, 0xFF)
			continue
		}
		checksum := crc8.Checksum([]byte(Normalize(linesSrc[line])), table)
		srcChk = append(srcChk, checksum)
	}
	return srcChk
}

// Normalize the line
// It will remove any character that is not a letter or a
// number included spaces, line feeds and tabs
func NormalizeOld(line string) string {

	var out string = ""

	for i := 0; i < len(line); i++ {
		if (line[i] >= '0' && line[i] <= '9') || (line[i] >= 'a' && line[i] <= 'z') {
			out += string(line[i])
		} else if line[i] >= 'A' && line[i] <= 'Z' {
			out += strings.ToLower(string(line[i]))
		}
	}
	return out

}

func Normalize(line string) string {

	var buffer bytes.Buffer
	//var out string = ""

	for i := 0; i < len(line); i++ {
		if (line[i] >= '0' && line[i] <= '9') || (line[i] >= 'a' && line[i] <= 'z') {
			buffer.WriteByte(line[i])
		} else if line[i] >= 'A' && line[i] <= 'Z' {
			buffer.WriteByte(line[i])
		}
	}
	return buffer.String()

}

// Get the longest snippet in the remote file that starts with a specific
// line in the local. A threshold must be reached to be considered as a match
//
func getSnippetsStarting(line uint32, localHashes []uint8, remoteHashes []uint8, remoteMap map[uint8][]uint32, Threshold uint32) (model.Range, int) {
	var snippet model.Range
	localStart := line
	remotes := remoteMap[localHashes[localStart]]
	l := 0
	err := 1
	for l = 0; l < len(remotes); {
		i := localStart
		j := remotes[l]
		for {
			if (int(i) < len(localHashes)) && (int(j) < len(remoteHashes)) && (localHashes[i] == remoteHashes[j]) {
				i++
				j++
			} else {
				break
			}
		}
		if (i - localStart) >= Threshold {
			if int(i-localStart) >= int(snippet.LEnd-snippet.LStart) {
				snippet.LStart = int(localStart)
				snippet.LEnd = int(i) - 1
				snippet.RStart = int(remotes[l])
				snippet.REnd = int(j) - 1
				err = 0

			}
		}
		l++
	}
	return snippet, err
}

// Compare search sequences of codes of local on the remote.
// A sequence is considered matched if at least reaches the Threshold
func Compare(local []uint8, remote []uint8, Threshold uint32) []model.Range {
	var ranges1 []model.Range
	var ranges2 []model.Range

	if Threshold == 0 {
		Threshold = 5
	}
	remoteMap := make(map[uint8][]uint32)
	var i int
	var j int
	exist := false
	for i = 0; i < len(remote); i++ {
		if hashes, ok := remoteMap[remote[i]]; ok {
			for j = range hashes {
				if hashes[j] == uint32(i) {
					exist = true
					break
				}
			}
			if !exist {
				hashes = append(hashes, uint32(i))
				remoteMap[remote[i]] = hashes
			}
		} else {
			hashes = append(hashes, uint32(i))
			remoteMap[remote[i]] = hashes
			//fmt.Printf("%x,%s", hash, purl)
		}

	}
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		for j := 0; j < len(local)/3; {
			a, err := getSnippetsStarting(uint32(j), local, remote, remoteMap, Threshold)
			if err == 0 {
				ranges1 = append(ranges1, a)
				j = a.LEnd + 1
			} else {
				j++
			}
		}
		wg.Done()
	}()
	go func() {
		for j := len(local)/2 + 1; j < len(local); {
			a, err := getSnippetsStarting(uint32(j), local, remote, remoteMap, Threshold)
			if err == 0 {
				ranges2 = append(ranges2, a)
				j = a.LEnd + 1
			} else {
				j++
			}
		}
		wg.Done()
	}()
	wg.Wait()
	finalRange := []model.Range{}
	l1 := len(ranges1)
	l2 := len(ranges2)

	if l1 > 0 && l2 > 0 {
		if ranges1[l1-1].REnd == ranges2[0].REnd {

			finalRange = append(ranges1, ranges2[1:]...)
		} else {
			finalRange = append(ranges1, ranges2...)
		}

	}
	/*
		for j = 0; j < len(local); {
			a, err := getSnippetsStarting(uint32(j), local, remote, remoteMap, Threshold)
			if err == 0 {
				ranges = append(ranges, a)
				j = a.LEnd + 1
			} else {
				j++
			}
		}*/
	//os.WriteFile("out.txt", []byte(fmt.Sprint(ranges1, "\n", ranges2, "\n", finalRange)), 0600)
	return finalRange
}
