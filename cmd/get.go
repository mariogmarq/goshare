package cmd

import (
	"io"
	"mime"
	"net/http"
	"os"

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

func get(cmd *cobra.Command, args []string) {
	//Search ip
	ip, err := util.ScanNetwork(":8080")
	if err == nil {
		resp, err := http.Get("http://" + ip + ":8080" + "/" + args[0])
		if err != nil {
			panic(err)
		}
		DownloadFile(resp)
	}
}

func DownloadFile(resp *http.Response) {
	defer resp.Body.Close()

	//Get the filename
	_, params, err := mime.ParseMediaType(resp.Header.Get("Content-Disposition"))
	if err != nil {
		panic(err)
	}

	file, err := os.Create(params["filename"])
	if err != nil {
		panic(err.Error())
	}

	io.Copy(file, resp.Body)

}
