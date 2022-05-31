package main

// use in your .go code

import (
	"fmt"

	proc "scanoss.com/hpsm/pkg"
)

func main() {
	hashLocal := proc.GetLineHashes("modified.txt")
	//hashRemote := proc.GetLineHashes("remote.c")
	/*fmt.Println(hashLocal)
	fmt.Println(hashRemote)*/

	//proc.Compare(hashLocal, hashRemote)
	for i := range hashLocal {
		fmt.Printf("%d, ", hashLocal[i])
	}

}

func Distance(local uint32, remote uint32) uint32 {
	if local == 0xFFFFFFFF && remote == 0xFFFFFFFF {
		return 32
	}
	aux := local ^ remote
	/*sum := aux & 0x000000FF
	aux = aux >> 8
	sum = aux & 0x000000FF
	aux = aux >> 8
	sum = aux & 0x000000FF
	aux = aux >> 8
	sum = aux & 0x000000FF
	*/
	var dist uint32 = 0
	for i := 0; i < 32; i++ {
		if (aux & 0x00000001) == 0x00000001 {
			dist++
		}
		aux = aux >> 1
	}
	return dist
}
