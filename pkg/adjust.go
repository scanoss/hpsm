package process

import (
	"os"
	"strings"

	"github.com/sigurn/crc8"
	model "scanoss.com/hpsm/model"
)

func Check(e error) {
	if e != nil {
		panic(e)
	}
}

func GetLineHashes(fileName string) []uint8 {
	src, err := os.ReadFile(fileName)
	Check(err)
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

func Normalize(line string) string {

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
				//return snippet, 0
			}
		} else {
		}
		l++
	}
	return snippet, err
}

// Compare search sequences of codes of local on the remote.
// A sequence is considered matched if at least reaches the Threshold
func Compare(local []uint8, remote []uint8, Threshold uint32) []model.Range {
	var ranges []model.Range
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
	for j = 0; j < len(local); {
		a, err := getSnippetsStarting(uint32(j), local, remote, remoteMap, Threshold)
		if err == 0 {
			ranges = append(ranges, a)
			j = a.LEnd + 1
		} else {
			j++
		}
	}
	return ranges
}
