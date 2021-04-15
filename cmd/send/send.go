package send

import (
	"bufio"
	"io"
	"net"
	"os"
	"sync"

	"github.com/mariogmarq/goshare/util"
	"github.com/urfave/cli/v2"
)

const HOST = "localhost"
const PORT = "13468"

// Main routine for the send command
func Send(c *cli.Context) error {
	//First we need to open the file and create a buffered reader
	filename := c.Args().First()
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	reader := bufio.NewReader(file)

	// Now create a tcp server
	l, err := net.Listen("tcp", HOST+":"+PORT)
	if err != nil {
		return err
	}
	defer l.Close()

	con, err := l.Accept()
	if err != nil {
		return err
	}

	//Handle the connection
	err = handleCon(con, reader)
	if err != nil {
		return err
	}

	return nil
}

//Handle the connection
func handleCon(con net.Conn, reader *bufio.Reader) error {
	lock := sync.WaitGroup{}

	buffer := make([]byte, reader.Size())
	writer := bufio.NewWriter(con)
	channel := make(chan []byte)

	lock.Add(1)
	util.Write(writer, channel, &lock)

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
