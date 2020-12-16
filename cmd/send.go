package cmd

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	"github.com/gin-gonic/gin"
	"github.com/mariogmarq/goshare/encryption"
	"github.com/mariogmarq/goshare/util"
	"github.com/spf13/cobra"
)

var (
	randomString  string
	encryptedData []byte
)

type key []byte

var sendCmd = &cobra.Command{
	Use:   "send [files to send] ",
	Short: "Allow you to share files accross the local network",
	Long:  "Allow you to share files accross the local network",
	Args:  cobra.MinimumNArgs(1),
	Run:   send,
}

func send(cmd *cobra.Command, args []string) {
	s := spinner.New(spinner.CharSets[11], 100*time.Millisecond, spinner.WithWriter(os.Stderr))
	s.Suffix = " Encrypting data..."
	s.Start()
	//Generate encryption key
	k, err := encryption.MakeKey(32)
	if err != nil {
		panic(err.Error())
	}

	//Load files into memory, verifing they exists(only one file supported)
	file, err := os.Open(args[0])
	if err != nil {
		panic("Error opening file")
	}

	//Encrypt the file
	encryptedData, err = encryption.Encrypt(k, file)
	if err != nil {
		panic("Error encrypting file")
	}
	s.Stop()
	//Once the file is ready, rise the server
	randomString = util.CreateRandomString(6)
	fmt.Printf("Code for share: %s\n", randomString)

	//Establish gin to Release mode
	gin.SetMode(gin.ReleaseMode)

	//generate the getHttpHandler
	getHandler := getSendHttpHandler(args[0])

	//Create http to listen to port
	g := gin.New()
	g.MaxMultipartMemory = 8 << 20 //8MB
	g.GET("/:code/get", getHandler)
	g.GET("/:code", pingHandler)
	g.GET("/:code/stop", stopHandler)
	g.GET("/:code/key", key(k).keyHandler)
	g.Run(":49153")
}

//Returns the handler for sending the file, takes name of file has a parameter
func getSendHttpHandler(filename string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Param("code") != randomString {
			return
		}
		fmt.Printf("Got connection with %s\n", c.Request.RemoteAddr)

		//Creates spinner
		s := spinner.New(spinner.CharSets[11], 100*time.Millisecond, spinner.WithWriter(os.Stderr))
		s.Suffix = " Sending data..."
		s.Start()

		//Establish the file name
		parsedFilename := strings.Split(filename, "/")
		//Write header for filename
		c.Writer.Header().Set("content-disposition",
			fmt.Sprintf("attachment; filename=\"%s\"",
				parsedFilename[len(parsedFilename)-1]))
		//Write the file
		c.Writer.Write(encryptedData)
		s.Stop()
		fmt.Println("File sent!")
	}
}

//Root handler just pings
func pingHandler(c *gin.Context) {
	if c.Param("code") != randomString {
		return
	}
	c.Header("status", "200")
}

//Writes the key into the response
func (k key) keyHandler(c *gin.Context) {
	if c.Param("code") != randomString {
		return
	}

	hexkey := fmt.Sprintf("%x", []byte(k))
	c.JSON(200, gin.H{
		"key": hexkey,
	})
}

//Stop handler for stop executing
func stopHandler(c *gin.Context) {
	if c.Param("code") != randomString {
		return
	}
	os.Exit(0)
}
