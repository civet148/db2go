package main

import (
	"fmt"
	_ "github.com/civet148/db2go/mssql"
	_ "github.com/civet148/db2go/mysql"
	_ "github.com/civet148/db2go/postgres"
	"github.com/civet148/db2go/schema"
	"github.com/civet148/log"
	"github.com/civet148/sqlca/v2"
	"github.com/urfave/cli/v2"
	"os"
	"os/signal"
	"strings"
)

const (
	SSH_SCHEME   = "ssh://"
	Version      = "2.2.0"
	ProgrameName = "db2go"
)

var (
	BuildTime = "2022-03-30"
	GitCommit = "<N/A>"
)

const (
	CMD_NAME_RUN = "run"
)

const (
	CmdFlag_Url            = "url"
	CmdFlag_Output         = "out"
	CmdFlag_Database       = "db"
	CmdFlag_Tables         = "table"
	CmdFlag_Tags           = "tag"
	CmdFlag_Prefix         = "prefix"
	CmdFlag_Suffix         = "suffix"
	CmdFlag_Package        = "package"
	CmdFlag_Without        = "without"
	CmdFlag_ReadOnly       = "readonly"
	CmdFlag_Protobuf       = "proto"
	CmdFlag_SpecType       = "spec-type"
	CmdFlag_EnableDecimal  = "enable-decimal"
	CmdFlag_GogoOptions    = "gogo-options"
	CmdFlag_OneFile        = "one-file"
	CmdFlag_DAO            = "orm"
	CmdFlag_OmitEmpty      = "omitempty"
	CmdFlag_JsonProperties = "json-properties"
	CmdFlag_TinyintAsBool  = "tinyint-as-bool"
	CmdFlag_SSH            = "ssh"
)

func init() {
	log.SetLevel("info")
}

func grace() {
	//capture signal of Ctrl+C and gracefully exit
	sigChannel := make(chan os.Signal, 1)
	signal.Notify(sigChannel, os.Interrupt)
	go func() {
		for {
			select {
			case s := <-sigChannel:
				{
					if s != nil && s == os.Interrupt {
						fmt.Printf("Ctrl+C signal captured, program exiting...\n")
						close(sigChannel)
						os.Exit(0)
					}
				}
			}
		}
	}()
}

func main() {

	grace()

	app := &cli.App{
		Name:    ProgrameName,
		Usage:   "db2go [options] --url <DSN>",
		Version: fmt.Sprintf("v%s %s commit %s", Version, BuildTime, GitCommit),
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     CmdFlag_Url,
				Usage:    "data source name of database",
				Required: true,
			},
			&cli.StringFlag{
				Name:        CmdFlag_Output,
				Usage:       "output path",
				DefaultText: ".",
			},
			&cli.StringFlag{
				Name:  CmdFlag_Database,
				Usage: "database name to export",
			},
			&cli.StringFlag{
				Name:  CmdFlag_Tables,
				Usage: "database tables to export",
			},
			&cli.StringFlag{
				Name:  CmdFlag_Tags,
				Usage: "export tags for golang",
			},
			&cli.StringFlag{
				Name:  CmdFlag_Prefix,
				Usage: "filename prefix",
			},
			&cli.StringFlag{
				Name:  CmdFlag_Suffix,
				Usage: "filename suffix",
			},
			&cli.StringFlag{
				Name:  CmdFlag_Package,
				Usage: "package name",
			},
			&cli.StringFlag{
				Name:  CmdFlag_Without,
				Usage: "exclude columns split by colon",
			},
			&cli.StringFlag{
				Name:  CmdFlag_ReadOnly,
				Usage: "readonly columns split by colon",
			},
			&cli.StringFlag{
				Name:  CmdFlag_Protobuf,
				Usage: "export protobuf file",
			},
			&cli.StringFlag{
				Name:  CmdFlag_SpecType,
				Usage: "specify column as customized types, e.g 'user.detail=UserDetail, user.data=UserData'",
			},
			&cli.BoolFlag{
				Name:  CmdFlag_EnableDecimal,
				Usage: "decimal as sqlca.Decimal type",
			},
			&cli.StringFlag{
				Name:  CmdFlag_GogoOptions,
				Usage: "gogo proto options",
			},
			&cli.BoolFlag{
				Name:  CmdFlag_OneFile,
				Usage: "export to one file",
			},
			&cli.StringFlag{
				Name:  CmdFlag_DAO,
				Usage: "generate data access object file",
				Value: "dao",
			},
			&cli.BoolFlag{
				Name:  CmdFlag_OmitEmpty,
				Usage: "json omitempty",
			},
			&cli.StringFlag{
				Name:  CmdFlag_JsonProperties,
				Usage: "customized properties for json tag",
			},
			&cli.StringFlag{
				Name:  CmdFlag_TinyintAsBool,
				Usage: "convert tinyint columns redeclare as bool type",
			},
			&cli.StringFlag{
				Name:  CmdFlag_SSH,
				Usage: "ssh tunnel e.g ssh://root:123456@192.168.1.23:22",
			},
		},
		Action: func(ctx *cli.Context) error {

			return doAction(ctx)
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Errorf("exit in error %s", err)
		os.Exit(1)
		return
	}
}

func doAction(ctx *cli.Context) error {
	//var err error
	var cmd = &schema.Commander{}
	cmd.Prefix = ctx.String(CmdFlag_Prefix)
	cmd.Suffix = ctx.String(CmdFlag_Suffix)
	cmd.OutDir = ctx.String(CmdFlag_Output)
	cmd.ConnUrl = ctx.String(CmdFlag_Url)
	cmd.PackageName = ctx.String(CmdFlag_Package)
	cmd.Protobuf = ctx.Bool(CmdFlag_Protobuf)
	cmd.EnableDecimal = ctx.Bool(CmdFlag_EnableDecimal)
	cmd.Orm = ctx.Bool(CmdFlag_DAO)
	cmd.OmitEmpty = ctx.Bool(CmdFlag_OmitEmpty)
	cmd.SSH = ctx.String(CmdFlag_SSH)
	cmd.Database = ctx.String(CmdFlag_Database)
	cmd.OneFile = ctx.Bool(CmdFlag_OneFile)
	cmd.JsonProperties = ctx.String(CmdFlag_JsonProperties)
	cmd.ParseSpecTypes(ctx.String(CmdFlag_SpecType))

	if cmd.SSH != "" {
		if !strings.Contains(cmd.SSH, SSH_SCHEME) {
			cmd.SSH = SSH_SCHEME + cmd.SSH
		}
	}

	ui := sqlca.ParseUrl(cmd.ConnUrl)

	log.Infof("%+v", cmd.String())

	if cmd.Database == "" {
		//use default database
		cmd.Database = schema.GetDatabaseName(ui.Path)
	}

	if ctx.String(CmdFlag_Tables) != "" {
		cmd.Tables = schema.TrimSpaceSlice(strings.Split(ctx.String(CmdFlag_Tables), ","))
	}

	if ctx.String(CmdFlag_Without) != "" {
		cmd.Without = schema.TrimSpaceSlice(strings.Split(ctx.String(CmdFlag_Without), ","))
	}

	if ctx.String(CmdFlag_TinyintAsBool) != "" {
		cmd.TinyintAsBool = schema.TrimSpaceSlice(strings.Split(ctx.String(CmdFlag_TinyintAsBool), ","))
	}

	if cmd.Protobuf {
		gogoOpt := ctx.String(CmdFlag_GogoOptions)
		if gogoOpt != "" {
			cmd.GogoOptions = schema.TrimSpaceSlice(strings.Split(gogoOpt, ","))
			if len(cmd.GogoOptions) == 0 {
				cmd.GogoOptions = schema.TrimSpaceSlice(strings.Split(gogoOpt, ";"))
			}
		}
	}

	if ctx.String(CmdFlag_Tags) != "" {
		cmd.Tags = schema.TrimSpaceSlice(strings.Split(ctx.String(CmdFlag_Tags), ","))
	}
	if ctx.String(CmdFlag_ReadOnly) != "" {
		cmd.ReadOnly = schema.TrimSpaceSlice(strings.Split(ctx.String(CmdFlag_ReadOnly), ","))
	}

	cmd.Scheme = ui.Scheme
	cmd.Host = ui.Host
	cmd.User = ui.User
	cmd.Password = ui.Password
	if strings.TrimSpace(cmd.SSH) != "" {
		cmd.Engine = sqlca.NewEngine(cmd.ConnUrl, tunnelOption(cmd.SSH))
	} else {
		cmd.Engine = sqlca.NewEngine(cmd.ConnUrl)
	}
	return export(cmd, cmd.Engine)
}

func tunnelOption(strSSH string) *sqlca.Options {
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

func export(cmd *schema.Commander, e *sqlca.Engine) (err error) {
	e.Debug(true)
	exporter := schema.NewExporter(cmd, e)
	if exporter == nil {
		err = fmt.Errorf("new exporter error, nil object")
		log.Error(err.Error())
		return err
	}
	if cmd.Protobuf {
		if err = exporter.ExportProto(); err != nil {
			log.Errorf("export [%v] to protobuf file error [%v]", cmd.Scheme, err.Error())
			return err
		}
	} else {
		if err := exporter.ExportGo(); err != nil {
			log.Errorf("export [%v] to go file error [%v]", cmd.Scheme, err.Error())
			return err
		}
	}
	return nil
}
