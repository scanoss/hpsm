package main

/*
int go_multiply(int a, int b);

typedef int (*multiply_f)(int a, int b);
multiply_f multiply;


static inline int multiplyWithFp(int a, int b) {
    return multiply(a, b);

}
 struct metadata_t{
	int size;
	long int md5;
};

*/
import "C"
import (
	"fmt"
	"os/exec"
	"unsafe"

	proc "scanoss.com/hpsm/pkg"
	u "scanoss.com/hpsm/utils"
)

/**
go build -o libhpsm.so  -buildmode=c-shared libhpsm.go
gcc -v client.c -o client ./libhpsm.so
*/

func Go_handleData(data *C.uchar, length C.int) []byte {
	return C.GoBytes(unsafe.Pointer(data), C.int(length))
}

func GetFileContent(url string, filepath string) error {
	// run shell `wget URL -O filepath`
	fmt.Printf("downloading %s -> %s\n", url, filepath)
	cmd := exec.Command("wget", url, "-O", filepath)
	return cmd.Run()
}

//export ProcessHPSM
func ProcessHPSM(data *C.uchar, length C.int, md5 *C.char) {
	dataArray := Go_handleData(data, length)
	MD5 := C.GoString(md5)
	//Remote access
	GetFileContent("https://osskb.org/api/file_contents/"+MD5, "/tmp/"+MD5)
	hashRemote := proc.GetLineHashes("/tmp/" + MD5)
	u.Rm("/tmp/" + MD5)
	hashLocal := dataArray
	snippets := proc.Compare(hashLocal, hashRemote, uint32(5))
	fmt.Println(snippets)
}
func main() {}
