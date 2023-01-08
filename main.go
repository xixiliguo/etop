package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
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
				"			relate value: 1h ago, 3h4m5s ago\n" +
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
			Usage:   "e.g 1h, 30m10s. if `DURATION` is specified, overrides --end value",
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
			Name:    "select",
			Aliases: []string{"s"},
			Value:   "",
			Usage:   "select `FIELD` for filter",
		},
		&cli.StringFlag{
			Name:    "filter",
			Aliases: []string{"F"},
			Value:   "",
			Usage:   "apply `REGEX` to filter --select field",
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
)

func normalizeField(module string, config model.RenderConfig, fields []string) (normalizedFields []string, err error) {

	if len(fields) == 1 {
		switch module {
		case "network":
			if v, ok := model.DefaultNetStatFields[fields[0]]; ok {
				return v, nil
			}
		}
	}
	for _, f := range fields {
		if _, ok := config[f]; !ok {
			return nil, fmt.Errorf("%s is not available field", f)
		}
	}
	return fields, nil
}

func dumpCommand(c *cli.Context, module string, fields []string) error {
	local, err := store.NewLocalStore(
		store.WithSetDefault(c.String("path"), logger),
	)
	if err != nil {
		return err
	}
	sm, err := model.NewSysModel(local, logger)
	if err != nil {
		return err
	}

	begin, err := util.ConvertToTime(c.String("begin"))
	if err != nil {
		return err
	}
	end := int64(0)
	if duration := c.String("duration"); duration == "" {
		end, err = util.ConvertToTime(c.String("end"))
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

	fields, err = normalizeField(module, sm.Config[module], fields)
	if err != nil {
		return err
	}

	if c.Bool("raw") == true {
		sm.Config[module].SetRawData()
	}
	output := os.Stdout

	out := c.String("output")
	if out != "" {
		if output, err = os.OpenFile(out, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644); err != nil {
			return err
		}
	}

	re, err := regexp.Compile(c.String("filter"))
	if err != nil {
		return err
	}
	opt := model.DumpOption{
		Begin:          begin,
		End:            end,
		Module:         module,
		Output:         output,
		Format:         c.String("output-format"),
		Fields:         fields,
		SelectField:    c.String("select"),
		Filter:         re,
		SortField:      c.String("sort"),
		AscendingOrder: c.Bool("ascending-order"),
		Top:            c.Int("top"),
		DisableTitle:   c.Bool("disable-title"),
		RepeatTitle:    c.Int("repeat-title"),
		RawData:        c.Bool("raw"),
	}
	return sm.Dump(opt)
}

var logger *log.Logger

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
					if _, err := os.Stat(path); os.IsNotExist(err) {
						if err := os.Mkdir(path, 0755); err != nil {
							return err
						}
					}
					local, err := store.NewLocalStore(
						store.WithSetDefault(path, logger),
						store.WithWriteOnly(),
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
							"			relate value: 1h ago, 3h4m5s ago\n" +
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
					if c.Bool("stat") == true {
						local, err := store.NewLocalStore(
							store.WithSetDefault(path, logger),
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
					t := tui.NewTUI(logger)
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
					t := tui.NewTUI(logger)
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
							"			relate value: 1h ago, 3h4m5s ago\n" +
							"			absolute value: 2006-01-02 15:04, 01-02 15:04, 15:04\n" +
							"			 ",
						Required: true,
					},
					&cli.StringFlag{
						Name:     "end",
						Aliases:  []string{"e"},
						Value:    "",
						Usage:    "`TIME` indicate end point of dump, same format with --begin",
						Required: true,
					},
					&cli.StringFlag{
						Name:    "duration",
						Aliases: []string{"d"},
						Value:   "",
						Usage:   "e.g 1h, 30m10s. if `DURATION` is specified, overrides --end value",
					},
					&cli.StringFlag{
						Name:    "path",
						Aliases: []string{"p"},
						Value:   "/var/log/etop",
						Usage:   "geneate snapshot file from `PATH`",
					},
				},
				Action: func(c *cli.Context) error {
					begin, err := util.ConvertToTime(c.String("begin"))
					if err != nil {
						return err
					}
					end := int64(0)
					if duration := c.String("duration"); duration == "" {
						end, err = util.ConvertToTime(c.String("end"))
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
					fmt.Printf("start to generate snapshot file from %s to %s\n",
						time.Unix(begin, 0),
						time.Unix(end, 0))
					local, err := store.NewLocalStore(
						store.WithSetDefault(c.String("path"), logger),
					)
					if err != nil {
						return err
					}
					if snapshotFileName, err := local.Snapshot(begin, end); err != nil {
						return err
					} else {
						fmt.Printf("snapshot file was created at %s\n", snapshotFileName)
					}
					return nil
				},
			},
			{
				Name:  "dump",
				Usage: "Dump data into parseable format",
				Subcommands: []*cli.Command{
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
							}),
						Action: func(c *cli.Context) error {
							fs := model.DefaultProcessFields
							if f := c.StringSlice("fields"); len(f) != 0 {
								fs = f
							}
							return dumpCommand(c, "process", fs)
						},
					},
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
