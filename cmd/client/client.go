package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"

	m "scanoss.com/hpsm/API/go"
	hashing "scanoss.com/hpsm/pkg"
	u "scanoss.com/hpsm/utils"
)

// ResultJson struct defines the structure for the
// results of a scanning process on SCANOSS platform
//
type ResultJson struct {
	Key         string   `json:"key,omitempty,omitempty"`
	ID          string   `json:"id,omitempty"`
	Status      string   `json:"status,omitempty"`
	Lines       string   `json:"lines,omitempty"`
	OssLines    string   `json:"oss_lines,omitempty"`
	Snippets    []string `json:"snippets,omitempty"`
	Matched     string   `json:"matched,omitempty"`
	Purl        []string `json:"purl,omitempty"`
	Vendor      string   `json:"vendor,omitempty"`
	Component   string   `json:"component,omitempty"`
	Version     string   `json:"version,omitempty"`
	Latest      string   `json:"latest,omitempty"`
	URL         string   `json:"url,omitempty"`
	ReleaseDate string   `json:"release_date,omitempty"`
	File        string   `json:"file,omitempty"`
	URLHash     string   `json:"url_hash,omitempty"`
	FileHash    string   `json:"file_hash,omitempty"`
	SourceHash  string   `json:"source_hash,omitempty"`
	FileURL     string   `json:"file_url,omitempty"`
	Licenses    []struct {
		Name             string `json:"name,omitempty"`
		PatentHints      string `json:"patent_hints,omitempty"`
		Copyleft         string `json:"copyleft,omitempty"`
		ChecklistURL     string `json:"checklist_url,omitempty"`
		IncompatibleWith string `json:"incompatible_with,omitempty,omitempty"`
		OsadlUpdated     string `json:"osadl_updated,omitempty"`
		Source           string `json:"source,omitempty"`
		Text             string `json:"license_text,omitempty"`
	} `json:"licenses,omitempty"`
	Server struct {
		Version   string `json:"version,omitempty"`
		KbVersion struct {
			Monthly string `json:"monthly,omitempty"`
			Daily   string `json:"daily,omitempty"`
		} `json:"kb_version,omitempty"`
	} `json:"server,omitempty"`
}
type ProcessJob struct {
	Result []ResultJson
	key    string
}

func Wget(url string, filepath string) error {
	cmd := exec.Command("wget", url, "-O", filepath, "-T", "10")
	return cmd.Run()
}

func main() {

	if len(os.Args) < 2 {
		fmt.Println("hpms <filename>")
		os.Exit(1)
	}

	x := map[string][]ResultJson{}

	// Run scanoss CLI
	cmd := exec.Command("scanoss-py", "scan", os.Args[1], "-T", "20")
	aux, err := cmd.Output()
	if err != nil {
		// Unmarshall results
		err = json.Unmarshal(aux, &x)
		if err != nil {
			log.Println(err)
		}
		var req []m.HpsmReqItem
		// From the results, create a list of files to be HPSM
		for key, val := range x {
			var job ProcessJob
			job.Result = val
			job.key = key
			for j := range val {
				if val[j].ID == "snippet" {
					var item m.HpsmReqItem
					item.MD5 = val[j].FileHash
					item.Hashes = hashing.GetLineHashes(os.Args[1] + key)
					req = append(req, item)
				}
			}
		}
		// Create the HPSM Req JSON
		out, _ := json.Marshal(req)
		//Request HPSM via CURL
		hpsm := u.Curl_HPSM("http://ns3193417.ip-152-228-225.eu:8081", string(out))
		//return scan results + HPSM
		fmt.Printf("{\"results\":%s\n,\"HPSM\": %s}", string(aux), hpsm)
	} else {
		fmt.Println("scanoss-py not detected or have no permissions")
	}

}
