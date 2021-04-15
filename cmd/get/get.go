package get

import (
	"bufio"
	"io"
	"net"
	"os"
	"sync"

	"github.com/mariogmarq/goshare/util"
	"github.com/urfave/cli/v2"
)

// Main routine for the get command
func Get(c *cli.Context) error {
	// TODO: Get this done, this variable is just for mocking
	address := "localhost:13468"

	// Try to get connection
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return err
	}
	defer conn.Close()

	// Handle the connection
	err = handleCon(conn, os.Stdout)
	if err != nil {
		return err
	}

	return nil
}

// Writes what it recives from conn in writer
// Is a decision of the programer to out to be buffered or not
func handleCon(conn net.Conn, out io.Writer) error {
	lock := sync.WaitGroup{}

	reader := bufio.NewReader(conn)
	buffer := make([]byte, reader.Size())

	channel := make(chan []byte)

	//Handle writing
	lock.Add(1)
	go util.Write(out, channel, &lock)

	for {
		_, err := reader.Read(buffer)
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		channel <- buffer
	}

	close(channel)
	lock.Wait()

	return nil
}
