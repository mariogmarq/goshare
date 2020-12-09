package cmd

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
)

var randomString string

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
	fmt.Printf("Code for share: %s\n", randomString)
	var files filesToSend = args

	//Establish gin to Release mode
	gin.SetMode(gin.ReleaseMode)

	//Create http to listen to port
	g := gin.Default()
	g.MaxMultipartMemory = 8 << 20 //8MB
	g.GET("/:code", files.sendHttpHandler)
	g.Run()
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

//Handler function for send
func (f filesToSend) sendHttpHandler(c *gin.Context) {
	if c.Param("code") == randomString {

		// Set Headers and print connection
		log.Printf("Got connection with %s\n", c.Request.RemoteAddr)

		for _, filename := range f {
			//Establish the file name
			parsedFilename := strings.Split(filename, "/")
			//Send files
			c.FileAttachment(filename, parsedFilename[len(parsedFilename)-1])
		}
		os.Exit(0)
	}
}
