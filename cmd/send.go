package cmd

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

var randomString string
var sendQuit chan bool

type filesToSend []string

var sendCmd = &cobra.Command{
	Use:   "send [files to send] ",
	Short: "Allow you to share files accross the local network",
	Long:  "Allow you to share files accross the local network",
	Args:  cobra.MinimumNArgs(1),
	Run:   send,
}

func send(cmd *cobra.Command, args []string) {
	randomString = createRandomString()
	fmt.Println("Code for share: " + randomString)
	var files filesToSend = args

	//Create http to listen to port
	http.HandleFunc("/", files.sendHandler)
	http.ListenAndServe(":8080", nil)
}

//Create a random string meant to be used in the send command
func createRandomString() string {
	rand.Seed(time.Now().Unix())
	charSet := "abcdedfghijklmnopqrstuvwxyz"
	var output strings.Builder
	length := 6
	for i := 0; i < length; i++ {
		random := rand.Intn(len(charSet))
		randomChar := charSet[random]
		output.WriteString(string(randomChar))
	}

	return output.String()
}

//handler for the send
func (f filesToSend) sendHandler(w http.ResponseWriter, req *http.Request) {

	//Wait for the request
	if req.Method == "GET" && req.URL.RawQuery == randomString {
		w.WriteHeader(200)
		fmt.Println("Got connection with " + string(req.RemoteAddr))

		for _, filename := range f {
			fmt.Println("Sending " + filename)
			file, err := os.Open(filename)

			if err != nil {
				fmt.Println("Error sending " + filename)
				continue
			}

			io.Copy(w, file)
		}
		os.Exit(0)
	}
}
