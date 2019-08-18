package main

import (
	"os"

	"github.com/fsnotify/fsnotify"
	"github.com/urfave/cli"
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
	app.Action = start
	app.Run(os.Args)
}

func start(c *cli.Context) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		panic(err)
	}
	defer watcher.Close()

	file_name := c.String("config")
	conf, err := read(file_name)
	if err != nil {
		panic(err)
	}

	for _, f := range conf.WatchFiles {
		err = watcher.Add(f)
		if err != nil {
			panic(err)
		}
	}

	done := make(chan struct{})

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					close(done)
					return
				}

				if event.Op == fsnotify.Write {
					runner := newRunner(conf)

					go runner.Start()
				}
			}
		}
	}()

	<-done

	return nil
}
