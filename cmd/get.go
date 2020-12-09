package cmd

import (
	"io"
	"log"
	"mime"
	"net"
	"net/http"
	"os"
	"time"

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
	client := http.Client{Timeout: 10 * time.Millisecond}

	ip, ipnet, err := net.ParseCIDR("192.168.1.0/24")
	if err != nil {
		log.Fatal(err)
	}

	//For every ip in your network
	for ip := ip.Mask(ipnet.Mask); ipnet.Contains(ip); inc(ip) {
		resp, err := client.Get("http://" + ip.String() + ":8080/" + args[0])
		if err == nil {
			DownloadFile(resp)
			break
		}
	}
}

func inc(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
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
