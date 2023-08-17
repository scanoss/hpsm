package main

/**

This file wraps HPSM functions to be used from another programming language
creating a shared library (.so or .dll)
Exported functions:
- HPSM: Formats an input string into []byte and send it with a MD5 of the OSS file. Computing HPSM is done by a remote server
- HashFileContents: This auxiliar function is used from a different language to calculate hashes

To build the library:
go build -o libhpsm.so  -buildmode=c-shared libhpsm.go
In order to buld a C client with the library
gcc -v client.c -o client ./libhpsm.so
*/

/*
 struct ranges{
	char *local;
	char *remote;
	char *matched;
};
*/
import "C"
import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
	"unsafe"

	"google.golang.org/grpc"
	pb "scanoss.com/hpsm/API/grpc"
	proc "scanoss.com/hpsm/pkg"
)

// Go_HandleData converts a unsigned char [] C array to an array
// of GO bytes
func Go_handleData(data *C.uchar, length C.int) []byte {
	return C.GoBytes(unsafe.Pointer(data), C.int(length))
}

//Auxiliar function to calculate hashes from a source code
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

//Calls HPSM caclulation on a gRPC service. The fist parameter is a HPSM definition (hpsm=zzyyww...) and
// the second is the key of the file to be compared. Returns a struct interpreted by C containting Ranges and matched percentage
//export HPSM
func HPSM(data *C.char, md5 *C.char) C.struct_ranges {
	dataArray := C.GoString(data)
	var crcSource []byte

	for i := 0; i < len(dataArray)-2; i += 2 {
		thisString := dataArray[i : i+2]
		thisByte, err := strconv.ParseInt(thisString, 16, 9)
		if err == nil {
			crcSource = append(crcSource, byte(thisByte))
		}
	}

	MD5 := C.GoString(md5)
	serverAddress := os.Getenv("HPSM_URL")
	if serverAddress == "" {
		serverAddress = "51.255.68.110:51015"
	}

	// Configures a context with timeout
	timeout := 10 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	conn, err := grpc.DialContext(ctx, serverAddress, grpc.WithInsecure())

	if err != nil {
		log.Printf("Failed to connect: %v", err)
		return C.struct_ranges{}
	}
	defer conn.Close()

	client := pb.NewHPSMClient(conn)

	response, err := client.ProcessHashes(context.Background(), &pb.HPSMRequest{Data: crcSource, Md5: MD5})
	if err != nil {
		log.Printf("Failed to process: %v", err)
		return C.struct_ranges{}
	}
	var lines C.struct_ranges
	lines.local = ((*C.char)(C.CString(response.Local)))
	lines.remote = ((*C.char)(C.CString(response.Remote)))
	lines.matched = ((*C.char)(C.CString(response.Matched)))
	return lines
}

//export ProcessHPSM
func ProcessHPSM(data *C.uchar, length C.int, md5 *C.char) C.struct_ranges {
	hashes := Go_handleData(data, length)
	MD5 := C.GoString(md5)
	conn, err := grpc.Dial("168.119.136.95:51015", grpc.WithInsecure())
	if err != nil {
		log.Printf("Failed to connect: %v", err)
		return C.struct_ranges{}
	}
	defer conn.Close()

	client := pb.NewHPSMClient(conn)
	response, err := client.ProcessHashes(context.Background(), &pb.HPSMRequest{Data: hashes, Md5: MD5})
	if err != nil {
		log.Printf("Failed to process: %v", err)
		return C.struct_ranges{}
	}
	var lines C.struct_ranges
	lines.local = ((*C.char)(C.CString(response.Local)))
	lines.remote = ((*C.char)(C.CString(response.Remote)))
	return lines
}
func main() {}
