package cmd

import (
	"github.com/mariogmarq/goshare/cmd/get"
	"github.com/urfave/cli/v2"
)

func Get(c *cli.Context) error {
	return get.Get(c)
}
