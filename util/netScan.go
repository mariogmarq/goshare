package util

import (
	"errors"
	"log"
	"net"
	"net/http"
	"time"
)

//Search for an ip in your network with the specified port opened
func ScanNetwork(port string) (string, error) {
	client := http.Client{Timeout: 10 * time.Millisecond}

	ip, ipnet, err := net.ParseCIDR("192.168.1.0/24")
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
