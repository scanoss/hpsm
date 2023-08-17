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
// Calculates CRC-8 of each line contained on a file given its path
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

// Hashes a string
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
func Normalize(line string) string {

	var buffer bytes.Buffer
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
// line in the local file. A threshold must be reached to be considered as a match
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

// Search common sequences of hashes from two array of hashes (refered as local and remote).
// A sequence is considered matched if at least reaches the Threshold number of contiguous hashes
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

		}

	}
	for j := 0; j < len(local); {
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

// Search common sequences of hashes from two array of hashes (refered as local and remote) by chunks.
// Each chunk is assigned to a different worker. Merging the results consists on checking if the last range of a chunk
// overlaps with the first of the following result
// A sequence is considered matched if at least reaches the Threshold
func CompareThreaded(local []uint8, remote []uint8, Threshold uint32, workers int) []model.Range {

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
		}

	}
	ranges := make([][]model.Range, workers)
	var wg sync.WaitGroup
	wg.Add(workers)

	for w := 0; w < workers; w++ {
		s := w * len(local) / workers
		e := (w + 1) * len(local) / workers

		go func(w int, start int, end int) {
			for j := start; j < end; {
				a, err := getSnippetsStarting(uint32(j), local, remote, remoteMap, Threshold)
				if err == 0 {
					ranges[w] = append(ranges[w], a)
					j = a.LEnd + 1
				} else {
					j++
				}
			}
			wg.Done()
		}(w, s, e)
	}
	wg.Wait()

	finalRange := []model.Range{}
	for r := 0; r < workers; r++ {
		l1 := len(finalRange)
		current := ranges[r]
		l2 := len(current)

		if l1 == 0 {
			finalRange = append(finalRange, current...)
		} else {
			if l2 > 0 && finalRange[l1-1].REnd == current[0].REnd {

				finalRange = append(finalRange, current[1:]...)
			} else {
				finalRange = append(finalRange, current...)
			}
		}
	}

	return finalRange
}
