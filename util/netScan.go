package util

import (
	"errors"
	"log"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/jackpal/gateway"
)

//Search for an ip in your network with the specified port opened
func ScanNetwork(port string) (string, error) {
	client := http.Client{Timeout: 10 * time.Millisecond}

	//Get the gateway for your network
	ipOfGateway, err := gateway.DiscoverGateway()
	if err != nil {
		return "", err
	}

	//Build the ip in a string
	var ipBuilder strings.Builder
	gatewayAsString := ipOfGateway.String()
	for i := 0; i < len(gatewayAsString)-1; i++ {
		ipBuilder.WriteByte(gatewayAsString[i])
	}

	ip, ipnet, err := net.ParseCIDR(ipBuilder.String() + "0/24")
	if err != nil {
		log.Fatal(err)
	}
	for ip := ip.Mask(ipnet.Mask); ipnet.Contains(ip); inc(ip) {
		_, err := client.Get("http://" + ip.String() + port)

		if err == nil {
			return ip.String(), nil
		}

	}

	return "", errors.New("Cannot find ip")
}

func inc(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}
