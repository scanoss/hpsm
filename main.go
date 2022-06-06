package main

// use in your .go code

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	proc "scanoss.com/hpsm/pkg"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("scan <hash|wfp> <filename>")
		fmt.Println("Available command")
		fmt.Println("hash : get lines hashes in one line from file")
		fmt.Println("wfp  : Fingerprints the file and adds the hash line")
	}
	if os.Args[1] == "hash" {
		hashLocal := proc.GetLineHashes(os.Args[2])
		fmt.Print("hpsm=")
		for i := range hashLocal {
			fmt.Printf("%02x", hashLocal[i])
		}
		os.Exit(0)
	}
	if os.Args[1] == "wfp" {
		cmd := exec.Command("scanoss-py", "wfp", os.Args[2])
		aux, _ := cmd.Output()

		lines := strings.Split(string(aux), "\n")
		out := lines[0] + "\n"

		// Unmarshall results
		hashLocal := proc.GetLineHashes(os.Args[2])
		out += ("hpsm=")
		for i := range hashLocal {
			out += fmt.Sprintf("%02x", hashLocal[i])
		}
		out += "\n"
		for j := 1; j < len(lines); j++ {
			out += lines[j] + "\n"
		}
		fmt.Println(out)

		os.Exit(0)
	}

}
