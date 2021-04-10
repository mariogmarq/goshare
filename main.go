package main

import (
	"log"
	"os"

	"github.com/mariogmarq/goshare/cmd"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Commands: []*cli.Command{
			{
				Name:   "send",
				Usage:  "send a file to other computer",
				Action: cmd.Send,
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
