// client.go
package main

import (
	"context"
	"fmt"
	"os"

	"google.golang.org/grpc"

	pb "scanoss.com/hpsm/API/grpc"
	hpsm "scanoss.com/hpsm/pkg"
)

func main() {
	conn, err := grpc.Dial("168.119.136.95:51015", grpc.WithInsecure())
	if err != nil {
		fmt.Printf("Failed to connect: %v", err)
		return
	}
	defer conn.Close()

	client := pb.NewHPSMClient(conn)

	hashes := hpsm.GetLineHashes(os.Args[1])
	md5 := os.Args[2]

	if err != nil {
		fmt.Printf("Failed to read file: %v", err)
		return
	}

	response, err := client.ProcessHashes(context.Background(), &pb.HPSMRequest{Data: hashes, Md5: md5})
	if err != nil {
		fmt.Printf("Failed to process: %v", err)
		return
	}

	fmt.Printf("Local: %s\n", response.Local)
	fmt.Printf("Remote: %s\n", response.Remote)
	fmt.Printf("Matched: %s\n", response.Matched)

}
