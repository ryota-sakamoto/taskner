package main

import (
	"os"

	"github.com/urfave/cli"

	"github.com/ryota-sakamoto/taskner/internal/taskner"
)

func main() {
	app := cli.NewApp()
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "config, c",
			Value: "config.yml",
			Usage: "job config file path",
		},
	}
	app.Action = taskner.WatchStart
	app.Run(os.Args)
}
