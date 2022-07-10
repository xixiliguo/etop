package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/urfave/cli/v2"
	"github.com/xixiliguo/etop/store"
	"github.com/xixiliguo/etop/tui"
	"github.com/xixiliguo/etop/version"
)

var logger *log.Logger

func main() {
	app := &cli.App{
		Name:    "etop",
		Version: version.Version,
		Usage:   "A tool for monitoring linux system resource",
		Commands: []*cli.Command{
			{
				Name:  "record",
				Usage: "record system data into file",
				Flags: []cli.Flag{
					&cli.IntFlag{
						Name:    "interval",
						Aliases: []string{"i"},
						Value:   5,
						Usage:   "second for interval",
					},
					&cli.StringFlag{
						Name:    "path",
						Aliases: []string{"p"},
						Value:   "/var/log/etop",
						Usage:   "`directory` to store file",
					},
				},
				Action: func(c *cli.Context) error {
					internal := c.Int("interval")
					if internal <= 0 {
						fmt.Printf("interval shoud great than 0, but get %d\n", internal)
						os.Exit(1)
					}
					local, err := store.NewLocalStore(c.String("path"), logger)
					if err != nil {
						return err
					}
					if err := local.WriteLoop(time.Duration(internal) * time.Second); err != nil {
						return err
					}
					return nil
				},
			},
			{
				Name:  "report",
				Usage: "Read file (created by etop record) and display",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "input",
						Aliases: []string{"i"},
						Value:   "",
						Usage:   "input file name",
					},
				},
				Action: func(c *cli.Context) error {
					t := tui.NewTUI(logger)
					if err := t.Run(c.String("input")); err != nil {
						return err
					}
					return nil
				},
			},
			{
				Name:  "live",
				Usage: "live display",
				Flags: []cli.Flag{
					&cli.IntFlag{
						Name:    "interval",
						Aliases: []string{"i"},
						Value:   5,
						Usage:   "second for interval",
					},
				},
				Action: func(c *cli.Context) error {
					internal := c.Int("interval")
					if internal <= 0 {
						fmt.Printf("interval shoud great than 0, but get %d\n", internal)
						os.Exit(1)
					}
					t := tui.NewTUI(logger)
					if err := t.RunWithLive(time.Duration(internal) * time.Second); err != nil {
						return err
					}
					return nil
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		logger.Fatalf("%s", err)
	}
}

func init() {
	file, err := os.OpenFile("/var/log/etop/etop.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Default().Fatalf("%s", err)
	}
	pidPrefix := fmt.Sprintf("[%d] ", os.Getpid())
	logger = log.New(file, pidPrefix, log.LstdFlags|log.Lshortfile)
}
