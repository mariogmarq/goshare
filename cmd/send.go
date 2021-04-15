package cmd

import (
	"github.com/mariogmarq/goshare/cmd/send"
	"github.com/urfave/cli/v2"
)

func Send(c *cli.Context) error {
	return send.Send(c)
}
