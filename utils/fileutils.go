package utils

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"
)

func Wget(url string, filepath string) error {

	out, err := os.Create(filepath)

	if err != nil {
		return err
	}
	http.DefaultClient.Timeout = 120 * time.Second
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error downloading ", err)
		return err
	}

	defer resp.Body.Close()
	// Check server response
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("%d", resp.StatusCode)
	}

	// Writer the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {

		return fmt.Errorf("error copying")
	}
	out.Close()
	return nil
}

func Mkdir(path string) error {
	cmd := exec.Command("mkdir", "-p", path)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
func Rm(file string) error {
	cmd := exec.Command("rm", file)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func Unzip(path string, dest string) error {
	cmd := exec.Command("unzip", "-PSecret", "-n", path+"/master.zip", "-d", dest)
	fmt.Printf("Extacting %s\n", path)
	return cmd.Run()

}

func Mz_Decompress(path string) (string, error) {
	cmd := exec.Command("mz", "-x", path)
	Mkdir("/tmp/scanoss/")
	cmd.Dir = ("/tmp/scanoss/")
	fmt.Printf("Extacting %s\n", path)
	return "/tmp/scanoss/", cmd.Run()

}

func Clean_dir(path string) error {
	cmd := exec.Command("rm", "-r", path)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()

}
func Count_lines(path string) int {

	//	lineStr := fmt.Sprintf("find %s -type f -exec wc -l {} \\; | awk '{total += $1} END{print total}'", path)
	value := 0
	cmd := exec.Command("bash", "utils/count.sh", path)
	out, _ := cmd.Output()

	fmt.Sscanf(string(out), "%d\n", &value)

	return value
}
func Get_Files(path string) []string {
	var ret []string
	fileInfo, err := os.Stat(path)
	if err != nil {
		fmt.Println("Error on stat of ", path)
	}

	if fileInfo.IsDir() {

		files, _ := ioutil.ReadDir(path)
		for fi := range files {
			ret = append(ret, Get_Files(path+"/"+files[fi].Name())...)
		}

		// is a directory
	} else {
		ret = append(ret, path)

	}
	return ret
}

func Scan(path string) string {

	//	lineStr := fmt.Sprintf("find %s -type f -exec wc -l {} \\; | awk '{total += $1} END{print total}'", path)

	cmd := exec.Command("scanoss-py", "scan", path)
	out, _ := cmd.Output()
	return string(out)
}
func Curl_HPSM(url string, req string) string {

	reader := strings.NewReader(req)
	request, err := http.NewRequest("POST", url+"/v2/adjust", reader)
	request.Header.Add("accept", "application/json")
	request.Header.Add("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(request)
	_ = err
	if err == nil {

		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)
		return string(body)
	}
	// TODO: check err
	return ""
}

func RequestHPSM(url string, req string) []byte {

	reader := strings.NewReader(req)
	request, _ := http.NewRequest("POST", url+"/v2/adjust", reader)
	request.Header.Add("accept", "application/json")
	request.Header.Add("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(request)
	_ = err
	if err == nil {

		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)
		return body
	}
	// TODO: check err
	return nil
}
