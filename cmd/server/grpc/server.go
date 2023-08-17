// server.go
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"

	"google.golang.org/grpc"
	pb "scanoss.com/hpsm/API/grpc"
	proc "scanoss.com/hpsm/pkg"
)

type HPSMServiceConfig struct {
	Port      int   `json:"port,omitempty"`
	Workers   int   `json:"workers,omitempty"`
	Threshold int32 `json:"threshold,omitempty"`
}
type server struct {
	pb.HPSMServer
}

var conf HPSMServiceConfig

func (s *server) ProcessHashes(ctx context.Context, req *pb.HPSMRequest) (*pb.RangeResponse, error) {

	strLinesGo := ""
	ossLinesGo := ""

	hashRemote := GetMD5Hashes(req.Md5)
	hashLocal := req.Data

	//snippets := proc.Compare(hashLocal, hashRemote, uint32(conf.Threshold))
	snippets := proc.CompareThreaded(hashLocal, hashRemote, uint32(conf.Threshold), conf.Workers)
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

	portStr := ":51015"
	if len(os.Args) == 2 {
		f, err := os.ReadFile(os.Args[1])
		if err != nil {
			log.Fatal("Could not open configuration file")
		}
		err = json.Unmarshal(f, &conf)
		if err != nil {
			log.Fatal("Could not load configuration parameters")
		}
		portStr = fmt.Sprintf(":%d", conf.Port)
	} else {
		conf.Threshold = 5
		conf.Workers = 2
	}
	listen, err := net.Listen("tcp", portStr)
	if err != nil {
		log.Printf("Failed to listen: %v", err)
		return
	}

	s := grpc.NewServer()
	pb.RegisterHPSMServer(s, &server{})

	log.Println("Server is listening on port 51015")
	if err := s.Serve(listen); err != nil {
		log.Printf("Failed to serve: %v", err)
	}
}
