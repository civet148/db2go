package main

import (
	"fmt"
	"os"
	"os/signal"
	"strings"

	_ "github.com/civet148/db2go/mssql"
	_ "github.com/civet148/db2go/mysql"
	_ "github.com/civet148/db2go/opengauss"
	_ "github.com/civet148/db2go/postgres"
	"github.com/civet148/db2go/schema"
	"github.com/civet148/log"
	"github.com/civet148/sqlca/v2"
	"github.com/urfave/cli/v2"
)

const (
	Version     = "v3.10.3"
	ProgramName = "db2go"
)

var (
	BuildTime = "2026-08-18"
	GitCommit = "<N/A>"
)

func init() {
	log.SetLevel("info")
}

func grace() {
	sigChannel := make(chan os.Signal, 1)
	signal.Notify(sigChannel, os.Interrupt)
	go func() {
		for {
			select {
			case s := <-sigChannel:
				if s != nil && s == os.Interrupt {
					fmt.Printf("Ctrl+C signal captured, program exiting...\n")
					close(sigChannel)
					os.Exit(0)
				}
			}
		}
	}()
}

func commonFlags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:     "url",
			Aliases:  []string{"u"},
			Usage:    "data source name of database",
			Required: true,
		},
		&cli.StringFlag{
			Name:    "out",
			Aliases: []string{"o"},
			Usage:   "output directory",
			Value:   ".",
		},
		&cli.StringFlag{
			Name:    "table",
			Aliases: []string{"t"},
			Usage:   "database tables to export (prefix with - to exclude)",
		},
		&cli.StringFlag{
			Name:    "package",
			Aliases: []string{"P"},
			Usage:   "package name",
		},
		&cli.StringFlag{
			Name:  "without",
			Usage: "exclude columns split by comma",
		},
		&cli.StringFlag{
			Name:  "ssh",
			Usage: "ssh tunnel e.g ssh://root:123456@192.168.1.23:22",
		},
		&cli.BoolFlag{
			Name:    "debug",
			Aliases: []string{"d"},
			Usage:   "open debug mode",
		},
	}
}

func goFlags() []cli.Flag {
	return append(commonFlags(),
		&cli.StringFlag{
			Name:    "tag",
			Aliases: []string{"T"},
			Usage:   "export tags for golang",
		},
		&cli.StringFlag{
			Name:    "readonly",
			Aliases: []string{"R"},
			Usage:   "readonly columns split by comma",
		},
		&cli.StringFlag{
			Name:    "spec-type",
			Aliases: []string{"S"},
			Usage:   "specify column as customized types, e.g 'user.detail=UserDetail, user.data=UserData'",
		},
		&cli.BoolFlag{
			Name:    "enable-decimal",
			Aliases: []string{"D"},
			Usage:   "decimal as sqlca.Decimal type",
		},
		&cli.StringFlag{
			Name:  "ddl",
			Usage: "export database DDL to file",
		},
	)
}

func protoFlags() []cli.Flag {
	return append(commonFlags(),
		&cli.StringFlag{
			Name:    "proto-options",
			Aliases: []string{"po"},
			Usage:   "set protobuf options, multiple options separated by ';'",
		},
		&cli.StringFlag{
			Name:    "gogo-options",
			Aliases: []string{"gogo"},
			Usage:   "gogo proto options",
		},
		&cli.StringFlag{
			Name:    "field-style",
			Aliases: []string{"style"},
			Usage:   "protobuf message field camel style (small or big)",
		},
	)
}

func main() {
	grace()

	app := &cli.App{
		Name:    ProgramName,
		Usage:   "db2go [command] [options]",
		Version: fmt.Sprintf("%s %s commit %s", Version, BuildTime, GitCommit),
		Commands: []*cli.Command{
			{
				Name:  "go",
				Usage: "Generate Go model files from database tables",
				Flags: goFlags(),
				Action: func(ctx *cli.Context) error {
					return runGo(ctx)
				},
			},
			{
				Name:  "proto",
				Usage: "Generate protobuf files from database tables",
				Flags: protoFlags(),
				Action: func(ctx *cli.Context) error {
					return runProto(ctx)
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		os.Exit(1)
	}
}

func buildCommonOptions(ctx *cli.Context) []schema.Option {
	var opts []schema.Option

	opts = append(opts, schema.WithURL(ctx.String("url")))
	opts = append(opts, schema.WithOutput(ctx.String("out")))
	opts = append(opts, schema.WithDatabase(ctx.String("db")))
	opts = append(opts, schema.WithTableFilter(ctx.String("table")))
	opts = append(opts, schema.WithPackageName(ctx.String("package")))
	opts = append(opts, schema.WithDebug(ctx.Bool("debug")))

	if v := ctx.String("without"); v != "" {
		opts = append(opts, schema.WithWithout(schema.TrimSpaceSlice(schema.Split(v))))
	}
	if v := ctx.String("spec-type"); v != "" {
		opts = append(opts, schema.WithSpecType(v))
	}
	opts = append(opts, schema.WithEnableDecimal(ctx.Bool("enable-decimal")))
	opts = append(opts, schema.WithSSH(ctx.String("ssh")))
	opts = append(opts, schema.WithFieldStyle(ctx.String("field-style")))
	opts = append(opts, schema.WithExportDDL(ctx.String("ddl")))

	return opts
}

func runGo(ctx *cli.Context) error {
	opts := buildCommonOptions(ctx)

	if v := ctx.String("tag"); v != "" {
		tags := schema.TrimSpaceSlice(schema.Split(v))
		opts = append(opts, schema.WithTags(tags))
	}
	if v := ctx.String("readonly"); v != "" {
		opts = append(opts, schema.WithReadOnly(schema.TrimSpaceSlice(schema.Split(v))))
	}

	cmd := schema.NewCmdFlags(opts...)
	log.Infof("开始生成golang文件...")
	return execute(cmd)
}

func runProto(ctx *cli.Context) error {
	opts := buildCommonOptions(ctx)

	if v := ctx.String("proto-options"); v != "" {
		protoOpts := make(map[string]string)
		for _, opt := range strings.Split(v, ";") {
			ss := strings.Split(opt, "=")
			if len(ss) != 2 {
				return fmt.Errorf("invalid protobuf option %s", opt)
			}
			protoOpts[ss[0]] = ss[1]
		}
		opts = append(opts, schema.WithProtoOptions(protoOpts))
	}
	if v := ctx.String("gogo-options"); v != "" {
		opts = append(opts, schema.WithGogoOptions(schema.TrimSpaceSlice(schema.Split(v))))
	}
	opts = append(opts, schema.WithProtobuf(true))

	cmd := schema.NewCmdFlags(opts...)
	log.Infof("开始生成protobuf文件...")
	return execute(cmd)
}

func execute(cmd *schema.CmdFlags) error {
	ui := sqlca.ParseUrl(cmd.ConnUrl)

	var opts = []schema.Option{
		schema.WithDatabase(schema.GetDatabaseName(ui.Path)),
		schema.WithScheme(ui.Scheme),
	}

	log.Json("command options", cmd)

	var err error
	var db *sqlca.Engine
	if strings.TrimSpace(cmd.SSH) != "" {
		db, err = sqlca.NewEngine(cmd.ConnUrl, sshOption(cmd.SSH))
	} else {
		db, err = sqlca.NewEngine(cmd.ConnUrl)
	}
	if err != nil {
		return log.Errorf("connect database [%s] error [%s]", cmd.ConnUrl, err.Error())
	}
	opts = append(opts, schema.WithEngine(db))

	for _, op := range opts {
		op(cmd)
	}

	exporter := schema.NewExporter(cmd, db)
	if exporter == nil {
		return log.Errorf("unsupported scheme: %s", cmd.Scheme)
	}

	if cmd.Protobuf {
		return exporter.ExportProto()
	}
	return exporter.ExportGo()
}

func sshOption(strSSH string) *sqlca.Options {
	if strSSH == "" {
		return nil
	}
	ssh := sqlca.ParseUrl(strSSH)
	return &sqlca.Options{
		SSH: &sqlca.SSH{
			User:     ssh.User,
			Password: ssh.Password,
			Host:     ssh.Host,
		},
	}
}
