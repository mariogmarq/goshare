package cmd

import (
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"mime"
	"net/http"
	"os"
	"time"

	"github.com/briandowns/spinner"
	"github.com/mariogmarq/goshare/encryption"
	"github.com/mariogmarq/goshare/util"
	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use:   "get [code for share] ",
	Short: "Allow you to get files accross the local network",
	Long:  "Allow you to get files accross the local network",
	Args:  cobra.MinimumNArgs(1),
	Run:   get,
}

var s *spinner.Spinner

func get(cmd *cobra.Command, args []string) {
	//Create spinner
	s = spinner.New(spinner.CharSets[11], 100*time.Millisecond, spinner.WithWriter(os.Stderr))
	s.Suffix = " Searching in network..."
	s.Start()

	//Search ip
	ip, err := util.ScanNetwork(":49153")
	if err == nil {
		s.Stop()
		ipHttp := "http://" + ip + ":49153/"
		//Get the key of encryption
		key, err := getKey(ipHttp + "key")
		if err != nil {
			panic(err.Error())
		}

		//Download the file
		resp, err := http.Get(ipHttp + "get/" + args[0])
		if err != nil {
			panic(err)
		}
		downloadFile(resp, key)

		//Shutdown the server
		http.Get(ipHttp + "stop")

	}
}

//Download the file and decrypts it
func downloadFile(resp *http.Response, key []byte) {
	defer resp.Body.Close()

	//Create spinner
	s.Suffix = " Donwloading and decrypting..."
	s.Start()

	//Get the filename
	_, params, err := mime.ParseMediaType(resp.Header.Get("Content-Disposition"))
	if err != nil {
		panic(err)
	}

	file, err := os.Create(params["filename"])
	if err != nil {
		panic(err.Error())
	}

	//Decrypt the file
	data, err := encryption.Decrypt(key, resp.Body)
	if err != nil {
		panic(err.Error())
	}

	file.Write(data)

	s.Stop()
}

//Get the key from the specified url
func getKey(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var key struct {
		Key string `json:"key"`
	}

	err = json.Unmarshal(data, &key)
	if err != nil {
		return nil, err
	}
	return hex.DecodeString(key.Key)
}
