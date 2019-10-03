package main

import (
	"github.com/bilfash/jingu/config"
	"github.com/bilfash/jingu/imap"
	"log"
	"os"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "jingu"
	app.Usage = "Download attachment from email"
	app.Version = "0.0.1"
	app.UsageText = "jingu path_to_config.yml"
	app.Action = func(c *cli.Context) error {
		configPath := c.Args().Get(0)
		config := config.New(configPath)
		err := imap.Download(config)
		return err
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
