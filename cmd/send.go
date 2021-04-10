package cmd

import (
	"bufio"
	"net"
	"os"

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
	buffer := make([]byte, reader.Size())
	writer := bufio.NewWriter(con)

	for {
		r, err := reader.Read(buffer)
		if err != nil {
			return err
		}

		if r == 0 {
			break
		}

		_, err = writer.Write(buffer)
		if err != nil {
			return err
		}

	}

	return nil
}
