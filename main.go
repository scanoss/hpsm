/*
 This file provides basic functionallity to demostrate HPSM in action
	It can be used as a CLI.
	Subcommands:
	- hash: Calculates HPSM sentence that should be included on a wfp file
	- compare: Compares two source codes or a source code and a remote MD5 in a splitted screen way
*/
package main

// use in your .go code

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	proc "scanoss.com/hpsm/pkg"
	"scanoss.com/hpsm/utils"
)

func setColor(c int) {
	colors := []string{"\033[31m ", "\033[32m", "\033[33m", "\033[34m", "\033[35m", "\033[36m", "\033[37m"}

	fmt.Println(string(colors[c]))
}
func gotoxy(x, y int) {
	fmt.Printf("\033[%d;%dH", x, y) // Set cursor position
}
func cls() {
	fmt.Print("\033[2J") //Clear screen
}

func trimLine(line string, maxLen int) string {
	if len(line) > maxLen {
		return line[:maxLen] + "...."
	} else {
		return line
	}
}

func main() {

	if len(os.Args) < 2 {
		fmt.Println("Available commands")
		fmt.Println("hash <filename>: gets line hashes in one line from file")
		fmt.Println("compare <localFile> <remoteFile|md5> [MD5]: Compares <localFile> against <remoteFile> or with remote <MD5>.")
		os.Exit(1)
	}
	if os.Args[1] == "hash" {
		hashLocal := proc.GetLineHashes(os.Args[2])
		fmt.Print("hpsm=")
		for i := range hashLocal {
			fmt.Printf("%02x", hashLocal[i])
		}
		os.Exit(0)
	}
	if os.Args[1] == "compare" {

		//setColor(2)
		var remote []byte
		var md5Int [2]uint
		matched, _ := fmt.Sscanf(os.Args[3], "%16x%16x", &md5Int[0], &md5Int[1])
		if matched == 2 {
			srcEndpoint := os.Getenv("SRC_URL")
			if srcEndpoint == "" {
				srcEndpoint = "https://osskb.org/api/file_contents/"
			}
			utils.Wget(srcEndpoint+os.Args[3], "/tmp/"+os.Args[3])
			remote, _ = os.ReadFile("/tmp/" + os.Args[3])
			utils.Rm("/tmp/" + os.Args[3])
		} else {
			remote, _ = os.ReadFile(os.Args[3])
		}

		src, _ := os.ReadFile(os.Args[2])
		linesSrc := strings.Split(string(src), "\n")

		linesRemote := strings.Split(string(remote), "\n")

		hashLocal := proc.GetLineHashes(os.Args[2])
		hashRemote := proc.GetLineHashesFromSource(string(remote))
		ranges := proc.CompareThreaded(hashLocal, hashRemote, 5, 10)
		y := 2

		for r := range ranges {
			cls()
			setColor(4)
			gotoxy(0, 10)
			fmt.Print("LOCAL SOURCE CODE")
			gotoxy(0, 90)
			fmt.Print("OSS SOURCE CODE")
			y = 2
			setColor(2)
			for l := -4; l < 0; l++ {
				gotoxy(y, 0)
				xL := l + ranges[r].LStart
				if xL > 0 && xL < len(linesSrc) {
					fmt.Print(xL, "\t", trimLine(linesSrc[xL], 30))
				} else {
					fmt.Print("\t", "[NO LINE]")
				}
				gotoxy(y, 80)
				xR := l + ranges[r].RStart
				if xR > 0 && xR < len(linesRemote)-1 {
					fmt.Print(xR, "\t", trimLine(linesRemote[xR], 30))
				} else {
					fmt.Print("\t", "[NO LINE]")
				}

				y++
			}
			setColor(3)
			for l := 0; l < (ranges[r].LEnd - ranges[r].LStart); l++ {
				gotoxy(y, 0)
				fmt.Print(l+ranges[r].LStart, "\t", trimLine(linesSrc[l+ranges[r].LStart], 30))
				gotoxy(y, 80)
				fmt.Print(l+ranges[r].RStart, "\t", trimLine(linesRemote[l+ranges[r].RStart], 30))
				fmt.Println()
				y++
			}
			setColor(2)
			for l := 0; l < 3; l++ {
				gotoxy(y, 0)
				xL := l + ranges[r].LEnd
				if xL > 0 && xL < len(linesSrc) {
					fmt.Print(xL, "\t", trimLine(linesSrc[xL], 30))
				} else {
					fmt.Print("[NO LINE]")
				}
				gotoxy(y, 80)
				xR := l + ranges[r].REnd
				if xR > 0 && xR < len(linesRemote)-1 {
					fmt.Print(xR, "\t", trimLine(linesRemote[xR], 30))
				} else {
					fmt.Print("[NO LINE]")
				}

				y++
			}
			reader := bufio.NewReader(os.Stdin)
			_, _ = reader.ReadString('\n')
			fmt.Println(("...."))
		}

		os.Exit(0)
	}

}
