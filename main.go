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
	SshScheme    = "ssh://"
	Version      = "2.7.0"
	ProgrameName = "db2go"
)

var (
	BuildTime = "2023-06-21"
	GitCommit = "<N/A>"
)

const (
	CmdFlag_Url            = "url"
	CmdFlag_Output         = "out"
	CmdFlag_Database       = "db"
	CmdFlag_DAO            = "dao"
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
	CmdFlag_OmitEmpty      = "omitempty"
	CmdFlag_JsonProperties = "json-properties"
	CmdFlag_TinyintAsBool  = "tinyint-as-bool"
	CmdFlag_SSH            = "ssh"
	CmdFlag_ImportModels   = "import-models"
	CmdFlag_V1             = "v1"
	CmdFlag_Debug          = "debug"
	CmdFlag_JsonStyle      = "json-style"
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
				Name:  CmdFlag_Output,
				Usage: "output path",
				Value: ".",
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
			&cli.BoolFlag{
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
			},
			&cli.StringFlag{
				Name:  CmdFlag_ImportModels,
				Usage: "project name",
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
			&cli.BoolFlag{
				Name:  CmdFlag_V1,
				Usage: "v1 package imports",
			},
			&cli.BoolFlag{
				Name:  CmdFlag_Debug,
				Usage: "open debug mode",
			},
			&cli.StringFlag{
				Name:  CmdFlag_JsonStyle,
				Usage: "json style: [underline or smallcamel/bigcamel] default underline",
				Value: schema.JSON_STYLE_UNDERLINE,
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
	cmd.Debug = ctx.Bool(CmdFlag_Debug)
	cmd.Prefix = ctx.String(CmdFlag_Prefix)
	cmd.Suffix = ctx.String(CmdFlag_Suffix)
	cmd.OutDir = ctx.String(CmdFlag_Output)
	cmd.ConnUrl = ctx.String(CmdFlag_Url)
	cmd.PackageName = ctx.String(CmdFlag_Package)
	cmd.Protobuf = ctx.Bool(CmdFlag_Protobuf)
	cmd.EnableDecimal = ctx.Bool(CmdFlag_EnableDecimal)
	cmd.DAO = ctx.String(CmdFlag_DAO)
	cmd.ImportModels = ctx.String(CmdFlag_ImportModels)
	cmd.OmitEmpty = ctx.Bool(CmdFlag_OmitEmpty)
	cmd.SSH = ctx.String(CmdFlag_SSH)
	cmd.Database = ctx.String(CmdFlag_Database)
	cmd.OneFile = ctx.Bool(CmdFlag_OneFile)
	cmd.JsonProperties = ctx.String(CmdFlag_JsonProperties)
	cmd.ParseSpecTypes(ctx.String(CmdFlag_SpecType))
	cmd.JsonStyle = ctx.String(CmdFlag_JsonStyle)

	if ctx.Bool(CmdFlag_V1) {
		cmd.ImportVer = schema.IMPORT_SQLCA_V1
	} else {
		cmd.ImportVer = schema.IMPORT_SQLCA_V2
	}
	if cmd.DAO != "" && cmd.ImportModels == "" {
		return log.Errorf("models path required eg. github.com/xxx/your-repo/models")
	}
	if cmd.SSH != "" {
		if !strings.Contains(cmd.SSH, SshScheme) {
			cmd.SSH = SshScheme + cmd.SSH
		}
	}

	if cmd.Debug {
		log.SetLevel("debug")
	}

	ui := sqlca.ParseUrl(cmd.ConnUrl)

	if cmd.Database == "" {
		//use default database
		cmd.Database = schema.GetDatabaseName(ui.Path)
		log.Infof("using default database %s", cmd.Database)
	}

	if ctx.String(CmdFlag_Tables) != "" {
		strFlagValue := ctx.String(CmdFlag_Tables)
		log.Infof("tables %+v", strFlagValue)
		tables := strings.Split(strFlagValue, ",")
		cmd.Tables = schema.TrimSpaceSlice(tables)
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

	log.Infof("command options [%+v]", cmd)

	var err error
	if strings.TrimSpace(cmd.SSH) != "" {
		cmd.Engine, err = sqlca.NewEngine(cmd.ConnUrl, Option(cmd.SSH))
	} else {
		cmd.Engine, err = sqlca.NewEngine(cmd.ConnUrl)
	}
	if err != nil {
		return log.Errorf("connect database [%s] error [%s]", cmd.ConnUrl, err.Error())
	}
	return export(cmd, cmd.Engine)
}

func Option(strSSH string) *sqlca.Options {
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
	exporter := schema.NewExporter(cmd, e)
	if exporter == nil {
		err = fmt.Errorf("new exporter error, nil object")
		log.Error(err.Error())
		return err
	}
	if cmd.Protobuf {
		log.Infof("generate protobuf files...")
		if err = exporter.ExportProto(); err != nil {
			log.Errorf("export [%v] to protobuf file error [%v]", cmd.Scheme, err.Error())
			return err
		}
	} else {
		log.Infof("generate golang files...")
		if err := exporter.ExportGo(); err != nil {
			log.Errorf("export [%v] to go file error [%v]", cmd.Scheme, err.Error())
			return err
		}
	}
	return nil
}
