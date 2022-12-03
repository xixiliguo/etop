package main

import (
	"fmt"
	"log"
	"os"
	"syscall"
	"time"

	"github.com/urfave/cli/v2"
	"github.com/xixiliguo/etop/store"
	"github.com/xixiliguo/etop/tui"
	"github.com/xixiliguo/etop/util"
	"github.com/xixiliguo/etop/version"
)

var (
	FreeSpaceForFileSystem uint64 = 500 * (1 << 20) //500MB
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
					intervalInt := c.Int("interval")
					if intervalInt <= 0 {
						fmt.Printf("interval shoud great than 0, but get %d\n", intervalInt)
						os.Exit(1)
					}
					local, err := store.NewLocalStore(
						store.WithSetDefault(c.String("path"), logger),
						store.WithWriteOnly(),
					)
					if err != nil {
						return err
					}
					// if err := local.WriteLoop(time.Duration(internal) * time.Second); err != nil {
					// 	return err
					// }
					interval := time.Duration(intervalInt) * time.Second
					logger.Printf("start to write sample every %s to %s", interval.String(), local.Data.Name())
					isSkip := 0
					for {
						start := time.Now()
						s := store.Sample{}
						if err := local.CollectSample(&s); err != nil {
							return err
						}

						statInfo := syscall.Statfs_t{}
						if err := syscall.Statfs(local.Path, &statInfo); err != nil {
							return err
						}
						if statInfo.Bavail*uint64(statInfo.Bsize) > FreeSpaceForFileSystem {
							if isSkip != 0 {
								logger.Printf("resume to write sample (%d skipped)", isSkip)
								isSkip = 0
							}
							if err := local.WriteSample(&s); err != nil {
								return err
							}
						} else {
							if isSkip == 0 {
								logger.Printf("filesystem free space %s below %s, write sample skipped",
									util.GetHumanSize(statInfo.Bavail*uint64(statInfo.Bsize)),
									util.GetHumanSize(FreeSpaceForFileSystem))
							}
							isSkip++
						}

						collectDuration := time.Now().Sub(start)
						if collectDuration > 500*time.Millisecond {
							logger.Printf("write sample take %s (larger than 500 ms)", collectDuration.String())
						}
						sleepDuration := time.Duration(1 * time.Second)
						if interval-collectDuration > 1*time.Second {
							sleepDuration = interval - collectDuration
						}
						time.Sleep(sleepDuration)
					}
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
