package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"github.com/urfave/cli/v2"
	"github.com/xixiliguo/etop/model"
	"github.com/xixiliguo/etop/store"
	"github.com/xixiliguo/etop/tui"
	"github.com/xixiliguo/etop/util"
	"github.com/xixiliguo/etop/version"
)

var (
	dumpFlag = []cli.Flag{
		&cli.StringFlag{
			Name:    "begin",
			Aliases: []string{"b"},
			Value:   time.Now().Format("2006-01-02") + " 00:00",
			Usage: "`TIME` indicate start point of dump\n" +
				"			relate value:  1h, mean 1 hour ago from now from now\n" +
				"			absolute value: 2006-01-02 15:04, 01-02 15:04, 15:04\n" +
				"			 ",
			DefaultText: time.Now().Format("2006-01-02") + " 00:00",
		},
		&cli.StringFlag{
			Name:        "end",
			Aliases:     []string{"e"},
			Value:       time.Now().Format("2006-01-02") + " 23:59",
			Usage:       "`TIME` indicate end point of dump, same format with --begin",
			DefaultText: time.Now().Format("2006-01-02") + " 23:59",
		},
		&cli.StringFlag{
			Name:    "duration",
			Aliases: []string{"d"},
			Value:   "",
			Usage:   "e.g 1h, 30m10s. if `DURATION` is specified, ignore --end value",
		},
		&cli.BoolFlag{
			Name:    "disable-title",
			Aliases: nil,
			Value:   false,
			Usage:   "disable title in text format output",
		},
		&cli.IntFlag{
			Name:    "repeat-title",
			Aliases: nil,
			Value:   0,
			Usage:   "for each `N` lines, it will render a line of title. only work for text format output",
		},
		&cli.StringSliceFlag{
			Name:    "fields",
			Aliases: []string{"f"},
			Value:   nil,
			Usage:   "select which `FIELD` to display",
		},
		&cli.StringFlag{
			Name:    "filter",
			Aliases: []string{"F"},
			Value:   "",
			Usage:   "apply filter",
		},
		&cli.StringFlag{
			Name:    "output",
			Aliases: []string{"o"},
			Value:   "",
			Usage:   "output destination, default to stdout",
		},
		&cli.StringFlag{
			Name:    "output-format",
			Aliases: []string{"O"},
			Value:   "text",
			Usage:   "output format, available value are text, json",
		},
		&cli.BoolFlag{
			Name:    "raw",
			Aliases: nil,
			Value:   false,
			Usage:   "dump raw data without units or conversion",
		},
		&cli.StringFlag{
			Name:    "path",
			Aliases: []string{"p"},
			Value:   "/var/log/etop",
			Usage:   "dump data from `PATH`",
		},
	}

	dumpOtelFlag = []cli.Flag{
		&cli.StringFlag{
			Name:    "begin",
			Aliases: []string{"b"},
			Value:   time.Now().Format("2006-01-02") + " 00:00",
			Usage: "`TIME` indicate start point of dump\n" +
				"			relate value:  1h, mean 1 hour ago from now from now\n" +
				"			absolute value: 2006-01-02 15:04, 01-02 15:04, 15:04\n" +
				"			 ",
			DefaultText: time.Now().Format("2006-01-02") + " 00:00",
		},
		&cli.StringFlag{
			Name:        "end",
			Aliases:     []string{"e"},
			Value:       time.Now().Format("2006-01-02") + " 23:59",
			Usage:       "`TIME` indicate end point of dump, same format with --begin",
			DefaultText: time.Now().Format("2006-01-02") + " 23:59",
		},
		&cli.StringFlag{
			Name:    "duration",
			Aliases: []string{"d"},
			Value:   "",
			Usage:   "e.g 1h, 30m10s. if `DURATION` is specified, ignore --end value",
		},

		&cli.StringFlag{
			Name:    "output",
			Aliases: []string{"o"},
			Value:   "",
			Usage:   "output destination, default to stdout",
		},

		&cli.StringFlag{
			Name:    "path",
			Aliases: []string{"p"},
			Value:   "/var/log/etop",
			Usage:   "dump data from `PATH`",
		},
	}
)

func dumpCommand(c *cli.Context, module string, fields []string) error {
	path := c.String("path")
	path, _ = filepath.Abs(path)
	logFile, err := os.OpenFile(filepath.Join(path, "etop.log"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	local, err := store.NewLocalStore(
		store.WithPathAndLogger(path, util.CreateLogger(logFile, false)),
	)
	if err != nil {
		return err
	}
	sm, err := model.NewSysModel(local, util.CreateLogger(logFile, false))
	if err != nil {
		return err
	}

	begin, err := util.ConvertToUnixTime(c.String("begin"))
	if err != nil {
		return err
	}
	end := int64(0)
	if duration := c.String("duration"); duration == "" {
		end, err = util.ConvertToUnixTime(c.String("end"))
		if err != nil {
			return err
		}
	} else {
		d, err := time.ParseDuration(duration)
		if err != nil {
			return err
		}
		end = begin + int64(d/time.Second)
	}

	output := os.Stdout

	out := c.String("output")
	if out != "" {
		if output, err = os.OpenFile(out, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644); err != nil {
			return err
		}
	}

	opt := model.DumpOption{
		Begin:           begin,
		End:             end,
		Module:          module,
		Output:          output,
		Format:          c.String("output-format"),
		Fields:          fields,
		FilterText:      c.String("filter"),
		SortField:       c.String("sort"),
		DescendingOrder: !c.Bool("ascending-order"),
		Top:             c.Int("top"),
		DisableTitle:    c.Bool("disable-title"),
		RepeatTitle:     c.Int("repeat-title"),
		RawData:         c.Bool("raw"),
	}
	return sm.Dump(opt)
}

func dumpToOtel(c *cli.Context) error {
	path := c.String("path")
	path, _ = filepath.Abs(path)
	logFile, err := os.OpenFile(filepath.Join(path, "etop.log"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	local, err := store.NewLocalStore(
		store.WithPathAndLogger(path, util.CreateLogger(logFile, false)),
	)
	if err != nil {
		return err
	}
	sm, err := model.NewSysModel(local, util.CreateLogger(logFile, false))
	if err != nil {
		return err
	}

	begin, err := util.ConvertToUnixTime(c.String("begin"))
	if err != nil {
		return err
	}
	end := int64(0)
	if duration := c.String("duration"); duration == "" {
		end, err = util.ConvertToUnixTime(c.String("end"))
		if err != nil {
			return err
		}
	} else {
		d, err := time.ParseDuration(duration)
		if err != nil {
			return err
		}
		end = begin + int64(d/time.Second)
	}

	opt := model.DumpOtelOption{
		Begin:  begin,
		End:    end,
		Output: c.String("output"),
	}

	return sm.DumpToOtel(opt)
}

func main() {
	app := &cli.App{
		Name:    "etop",
		Version: version.Version,
		Usage:   "A tool for monitoring linux system resource",
		Commands: []*cli.Command{
			{
				Name:  "record",
				Usage: "Record system resource into file",
				Flags: []cli.Flag{
					&cli.IntFlag{
						Name:    "interval",
						Aliases: []string{"i"},
						Value:   5,
						Usage:   "number of seconds between samples",
					},
					&cli.StringFlag{
						Name:    "path",
						Aliases: []string{"p"},
						Value:   "/var/log/etop",
						Usage:   "`directory` to store file",
					},
					&cli.BoolFlag{
						Name:  "compress",
						Value: false,
						Usage: "enable zstd compress",
					},
					&cli.IntFlag{
						Name:  "compress-dict-chunk-size",
						Value: 0,
						Usage: "at least 2, must be a power of two",
						Action: func(c *cli.Context, n int) error {
							if c.Bool("compress") == false {
								return fmt.Errorf("compress-dict-chunk-size only is valid when --compress")
							}
							if n > store.MaxDictOffset {
								return fmt.Errorf("value must <= %d", store.MaxDictOffset)
							}
							if n >= 2 && ((n & (n - 1)) == 0) {
							} else {
								return fmt.Errorf("compress-dict-chunk-size at least 2, must be a power of two")
							}
							return nil
						},
					},
					&cli.IntFlag{
						Name:  "retainday",
						Value: 3,
						Usage: "maximum days for retaining data file, detele oldest one if exceed `THRESHOLD`",
					},
					&cli.Int64Flag{
						Name:        "retainsize",
						Value:       20 * (1 << 30), // 20GB,
						DefaultText: "20 GB",
						Usage:       "size limit in bytes for retaining data file, detele oldest one if exceed `THRESHOLD`",
					},
				},
				Action: func(c *cli.Context) error {
					intervalFlag := c.Int("interval")
					if intervalFlag <= 0 {
						return fmt.Errorf("interval flag shoud great than 0, but get %d\n", intervalFlag)
					}
					retaindayFlag := c.Int("retainday")
					if retaindayFlag <= 0 {
						return fmt.Errorf("retainday flag shoud great than 0, but get %d\n", retaindayFlag)
					}
					retainsizeFlag := c.Int64("retainsize")
					if retainsizeFlag <= 0 {
						return fmt.Errorf("retainday flag shoud great than 0, but get %d\n", retainsizeFlag)
					}
					path := c.String("path")
					path, _ = filepath.Abs(path)
					if _, err := os.Stat(path); os.IsNotExist(err) {
						if err := os.Mkdir(path, 0755); err != nil {
							return err
						}
					}

					logFile, err := os.OpenFile(filepath.Join(path, "etop.log"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
					if err != nil {
						return err
					}

					log := util.CreateLogger(logFile, false)
					msg := fmt.Sprintf("version: %s", version.Version)
					log.Info(msg)

					mode := store.NoCompress
					if c.Bool("compress") {
						mode = store.ZstdCompress
					}
					chunk := uint32(c.Int("compress-dict-chunk-size"))
					if chunk != 0 {
						mode = store.ZstdCompressWithDict
					}
					local, err := store.NewLocalStore(
						store.WithPathAndLogger(path, log),
						store.WithWriteOnly(mode, chunk),
						store.WithExitProcess(log),
					)
					if err != nil {
						return err
					}
					opt := store.WriteOption{
						Interval:   time.Duration(intervalFlag) * time.Second,
						RetainDay:  retaindayFlag,
						RetainSize: retainsizeFlag,
					}
					if err := local.WriteLoop(opt); err != nil {
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
						Name:    "path",
						Aliases: []string{"p"},
						Value:   "/var/log/etop",
						Usage:   "read data from `PATH`",
					},
					&cli.StringFlag{
						Name:    "begin",
						Aliases: []string{"b"},
						Value:   time.Now().Format("2006-01-02") + " 00:00",
						Usage: "`TIME` indicate start point of report\n" +
							"			relate value:  1h, mean 1 hour ago from now\n" +
							"			absolute value: 2006-01-02 15:04, 01-02 15:04, 15:04\n" +
							"			 ",
						DefaultText: time.Now().Format("2006-01-02") + " 00:00",
					},
					&cli.BoolFlag{
						Name:    "stat",
						Aliases: []string{"s"},
						Value:   false,
						Usage:   "output data statistics info instead of display",
					},
					&cli.StringFlag{
						Name:  "snapshot",
						Value: "",
						Usage: "read data from snapshot `FILE`",
					},
				},
				Action: func(c *cli.Context) error {
					path := c.String("path")
					if snapshot := c.String("snapshot"); snapshot != "" {
						if tempPath, err := util.ExtractFileFromTar(snapshot); err != nil {
							return err
						} else {
							path = tempPath
						}
					}
					path, _ = filepath.Abs(path)
					if c.Bool("stat") == true {

						local, err := store.NewLocalStore(
							store.WithPathAndLogger(path, util.CreateLogger(os.Stdout, false)),
						)
						if err != nil {
							return err
						}
						result, err := local.FileStatInfo()
						if err != nil {
							return err
						}
						fmt.Println(result)
						return nil
					}

					t := tui.NewTUI()
					if err := t.Run(path, c.String("begin")); err != nil {
						return err
					}
					return nil
				},
			},
			{
				Name:  "live",
				Usage: "Live display",
				Flags: []cli.Flag{
					&cli.IntFlag{
						Name:    "interval",
						Aliases: []string{"i"},
						Value:   5,
						Usage:   "number of seconds between samples",
					},
				},
				Action: func(c *cli.Context) error {
					internal := c.Int("interval")
					if internal <= 0 {
						fmt.Printf("interval shoud great than 0, but get %d\n", internal)
						os.Exit(1)
					}
					t := tui.NewTUI()
					if err := t.RunWithLive(time.Duration(internal) * time.Second); err != nil {
						return err
					}
					return nil
				},
			},
			{
				Name:  "snapshot",
				Usage: "Generate a snapshot file based on the time range provided",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "begin",
						Aliases: []string{"b"},
						Value:   "",
						Usage: "`TIME` indicate start point of dump\n" +
							"			relate value:  1h, mean 1 hour ago from now\n" +
							"			absolute value: 2006-01-02 15:04, 01-02 15:04, 15:04\n" +
							"			 ",
						Required: true,
						Action: func(c *cli.Context, v string) error {
							if c.String("end") == "" && c.String("duration") == "" {
								return fmt.Errorf("either of end, duration options must be specified")
							}
							return nil
						},
					},
					&cli.StringFlag{
						Name:    "end",
						Aliases: []string{"e"},
						Value:   "",
						Usage:   "`TIME` indicate end point of dump, same format with --begin",
					},
					&cli.StringFlag{
						Name:    "duration",
						Aliases: []string{"d"},
						Value:   "",
						Usage:   "e.g 1h, 30m10s. if `DURATION` is specified, ignore --end value",
					},
					&cli.StringFlag{
						Name:    "path",
						Aliases: []string{"p"},
						Value:   "/var/log/etop",
						Usage:   "geneate snapshot file from `PATH`",
					},
				},
				Action: func(c *cli.Context) error {
					begin, err := util.ConvertToUnixTime(c.String("begin"))
					if err != nil {
						return err
					}
					end := int64(0)
					if duration := c.String("duration"); duration == "" {
						end, err = util.ConvertToUnixTime(c.String("end"))
						if err != nil {
							return err
						}
					} else {
						d, err := time.ParseDuration(duration)
						if err != nil {
							return err
						}
						end = begin + int64(d/time.Second)
					}
					msg := fmt.Sprintf("start to generate snapshot file from %s to %s\n",
						time.Unix(begin, 0),
						time.Unix(end, 0))
					slog.Info(msg)
					local, err := store.NewLocalStore(
						store.WithPathAndLogger(c.String("path"), slog.Default()),
					)
					if err != nil {
						return err
					}
					if snapshotFileName, err := local.Snapshot(begin, end); err != nil {
						return err
					} else {
						msg := fmt.Sprintf("snapshot file was created at %s\n", snapshotFileName)
						slog.Info(msg)
					}
					return nil
				},
			},
			{
				Name:  "dump",
				Usage: "Dump data into parseable format",
				Subcommands: []*cli.Command{
					{
						Name:  "system",
						Usage: "Dump system stat",
						Flags: dumpFlag,
						Action: func(c *cli.Context) error {
							fs := model.DefaultSystemFields
							if f := c.StringSlice("fields"); len(f) != 0 {
								fs = f
							}
							return dumpCommand(c, "system", fs)
						},
					},
					{
						Name:  "cpu",
						Usage: "Dump cpu stat",
						Flags: dumpFlag,
						Action: func(c *cli.Context) error {
							fs := model.DefaultCPUFields
							if f := c.StringSlice("fields"); len(f) != 0 {
								fs = f
							}
							return dumpCommand(c, "cpu", fs)
						},
					},
					{
						Name:  "memory",
						Usage: "Dump memory stat",
						Flags: dumpFlag,
						Action: func(c *cli.Context) error {
							fs := model.DefaultMEMFields
							if f := c.StringSlice("fields"); len(f) != 0 {
								fs = f
							}
							return dumpCommand(c, "memory", fs)
						},
					},
					{
						Name:  "vm",
						Usage: "Dump vm stat",
						Flags: dumpFlag,
						Action: func(c *cli.Context) error {
							fs := model.DefaultVmFields
							if f := c.StringSlice("fields"); len(f) != 0 {
								fs = f
							}
							return dumpCommand(c, "vm", fs)
						},
					},
					{
						Name:  "disk",
						Usage: "Dump disk stat",
						Flags: dumpFlag,
						Action: func(c *cli.Context) error {
							fs := model.DefaultDiskFields
							if f := c.StringSlice("fields"); len(f) != 0 {
								fs = f
							}
							return dumpCommand(c, "disk", fs)
						},
					},
					{
						Name:  "netdev",
						Usage: "Dump netdev stat",
						Flags: dumpFlag,
						Action: func(c *cli.Context) error {
							fs := model.DefaultNetDevFields
							if f := c.StringSlice("fields"); len(f) != 0 {
								fs = f
							}
							return dumpCommand(c, "netdev", fs)
						},
					},
					{
						Name:  "network",
						Usage: "Dump network snmp and ext stat",
						Flags: dumpFlag,
						Action: func(c *cli.Context) error {
							fs := model.DefaultNetStatFields["tcp"]
							if f := c.StringSlice("fields"); len(f) != 0 {
								fs = f
							}
							return dumpCommand(c, "network", fs)
						},
					},
					{
						Name:  "networkprotocol",
						Usage: "Dump networkprotocol stat",
						Flags: dumpFlag,
						Action: func(c *cli.Context) error {
							fs := model.DefaultNetProtocolFields
							if f := c.StringSlice("fields"); len(f) != 0 {
								fs = f
							}
							return dumpCommand(c, "networkprotocol", fs)
						},
					},
					{
						Name:  "softnet",
						Usage: "Dump softnet stat",
						Flags: dumpFlag,
						Action: func(c *cli.Context) error {
							fs := model.DefaultSoftnetFields
							if f := c.StringSlice("fields"); len(f) != 0 {
								fs = f
							}
							return dumpCommand(c, "softnet", fs)
						},
					},
					{
						Name:  "process",
						Usage: "Dump process stat",
						Flags: append(dumpFlag,
							&cli.StringFlag{
								Name:  "sort",
								Value: "CPU",
								Usage: "sort `FIELD` by descending order",
							},
							&cli.BoolFlag{
								Name:    "ascending-order",
								Aliases: nil,
								Value:   false,
								Usage:   "sort by ascending order",
							},
							&cli.IntFlag{
								Name:  "top",
								Value: 0,
								Usage: "show top `N` info",
							},
							&cli.BoolFlag{
								Name:  "all",
								Value: false,
								Usage: "dump all fields",
							}),
						Action: func(c *cli.Context) error {
							fs := model.DefaultProcessFields
							if c.Bool("all") == true {
								fs = model.AllProcessFields
							}
							if f := c.StringSlice("fields"); len(f) != 0 {
								fs = f
							}
							return dumpCommand(c, "process", fs)
						},
					},
					{
						Name:  "cgroup",
						Usage: "Dump cgroup stat",
						Flags: dumpFlag,
						Action: func(c *cli.Context) error {
							fs := model.DefaultCgroupFields
							if c.Bool("all") == true {
								fs = model.AllCgroupFields
							}
							if f := c.StringSlice("fields"); len(f) != 0 {
								fs = f
							}
							return dumpCommand(c, "cgroup", fs)
						},
					},
					{
						Name:  "otel",
						Usage: "Dump and send to otel backend",
						Flags: dumpOtelFlag,
						Action: func(c *cli.Context) error {
							return dumpToOtel(c)
						},
					},
				},
			},
			{
				Name:  "debug",
				Usage: "Provides various debug feature",
				Subcommands: []*cli.Command{
					{
						Name:  "dump-store",
						Usage: "Dump raw data of store",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:    "begin",
								Aliases: []string{"b"},
								Value:   "",
								Usage: "`TIME` indicate start point of dump\n" +
									"			relate value:  1h, mean 1 hour ago from now\n" +
									"			absolute value: 2006-01-02 15:04, 01-02 15:04, 15:04\n" +
									"			 ",
								Required: true,
								Action: func(c *cli.Context, v string) error {
									if c.String("end") == "" && c.String("duration") == "" {
										return fmt.Errorf("either of end, duration options must be specified")
									}
									return nil
								},
							},
							&cli.StringFlag{
								Name:    "end",
								Aliases: []string{"e"},
								Value:   "",
								Usage:   "`TIME` indicate end point of dump, same format with --begin",
							},
							&cli.StringFlag{
								Name:    "duration",
								Aliases: []string{"d"},
								Value:   "",
								Usage:   "e.g 1h, 30m10s. if `DURATION` is specified, ignore --end value",
							},
							&cli.StringFlag{
								Name:    "path",
								Aliases: []string{"p"},
								Value:   "/var/log/etop",
								Usage:   "geneate snapshot file from `PATH`",
							},
						},
						Action: func(c *cli.Context) error {
							path := c.String("path")
							path, _ = filepath.Abs(path)
							local, err := store.NewLocalStore(
								store.WithPathAndLogger(path, util.CreateLogger(os.Stdout, false)),
							)
							if err != nil {
								return err
							}
							begin, err := util.ConvertToUnixTime(c.String("begin"))
							if err != nil {
								return err
							}
							end := int64(0)
							if duration := c.String("duration"); duration == "" {
								end, err = util.ConvertToUnixTime(c.String("end"))
								if err != nil {
									return err
								}
							} else {
								d, err := time.ParseDuration(duration)
								if err != nil {
									return err
								}
								end = begin + int64(d/time.Second)
							}
							sample := store.NewSample()
							if err := local.JumpSampleByTimeStamp(begin, &sample); err != nil {
								return err
							}
							for begin < end {
								b, err := json.MarshalIndent(sample, "", "	")
								if err != nil {
									return err
								}
								fmt.Printf("%s\n", b)
								sample.Reset()
								if err := local.NextSample(1, &sample); err != nil {
									return nil
								}
								begin = sample.TimeStamp
							}
							return nil
						},
					},
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
