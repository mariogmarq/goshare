package cmd

import (
	"io"
	"log"
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
		req, err := client.Get("http://" + ip.String() + ":8080/" + args[0])
		if err == nil {
			defer req.Body.Close()
			file, err := os.Create(args[0])
			if err != nil {
				panic(err.Error())
			}

			io.Copy(file, req.Body)

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
