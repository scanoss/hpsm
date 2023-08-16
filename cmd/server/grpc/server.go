// server.go
package main

import (
	"context"
	"fmt"
	"net"
	"os/exec"

	"google.golang.org/grpc"

	pb "scanoss.com/hpsm/API/grpc" // Cambia "your_package_path" al directorio donde está tu archivo .proto generado
	proc "scanoss.com/hpsm/pkg"
)

type server struct {
	pb.HPSMServer
}

func (s *server) ProcessHashes(ctx context.Context, req *pb.HPSMRequest) (*pb.RangeResponse, error) {

	strLinesGo := ""
	ossLinesGo := ""

	hashRemote := GetMD5Hashes(req.Md5)
	hashLocal := req.Data

	snippets := proc.Compare(hashLocal, hashRemote, uint32(5))
	matchedLines := 0
	totalLines := len(hashLocal)
	for i := range snippets {
		var ossRange string
		var srcRange string

		matchedLines += (snippets[i].LEnd - snippets[i].LStart)
		srcRange = fmt.Sprintf("%d-%d", snippets[i].LStart, snippets[i].LEnd)
		ossRange = fmt.Sprintf("%d-%d", snippets[i].RStart, snippets[i].REnd)
		strLinesGo += ossRange
		ossLinesGo += srcRange
		if i < len(snippets)-1 {
			strLinesGo += ", "
			ossLinesGo += ", "
		}
	}
	mLines := ""
	if len(snippets) == 0 {
		mLines = "0%"
	} else {
		mLines = fmt.Sprintf("%d%%", matchedLines*100/totalLines)
	}

	return &pb.RangeResponse{Local: strLinesGo, Remote: ossLinesGo, Matched: mLines}, nil
}
func GetMD5Hashes(md5key string) []byte {

	cmd := exec.Command("scanoss", "-k", md5key)
	out, _ := cmd.Output()

	return proc.GetLineHashesFromSource(string(out))

}
func main() {
	listen, err := net.Listen("tcp", ":51015")
	if err != nil {
		fmt.Printf("Failed to listen: %v", err)
		return
	}

	s := grpc.NewServer()
	pb.RegisterHPSMServer(s, &server{})

	fmt.Println("Server is listening on port 51015")
	if err := s.Serve(listen); err != nil {
		fmt.Printf("Failed to serve: %v", err)
	}
}
