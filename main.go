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
	SshScheme   = "ssh://"
	Version     = "3.4.2"
	ProgramName = "db2go"
)

var (
	BuildTime = "2025-12-31"
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
	CmdFlag_ProtoOptions   = "proto-options"
	CmdFlag_Merge          = "merge"
	CmdFlag_OmitEmpty      = "omitempty"
	CmdFlag_JsonProperties = "json-properties"
	CmdFlag_TinyintAsBool  = "tinyint-as-bool"
	CmdFlag_SSH            = "ssh"
	CmdFlag_ImportModels   = "import-models"
	CmdFlag_V2             = "v2"
	CmdFlag_Debug          = "debug"
	CmdFlag_JsonStyle      = "json-style"
	CmdFlag_Export         = "export"
	CmdFlag_FieldStyle     = "field-style"
	CmdFlag_BaseModel      = "base-model"
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
		Name:    ProgramName,
		Usage:   "db2go [options] --url <DSN>",
		Version: fmt.Sprintf("v%s %s commit %s", Version, BuildTime, GitCommit),
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     CmdFlag_Url,
				Aliases:  []string{"u"},
				Usage:    "data source name of database",
				Required: true,
			},
			&cli.StringFlag{
				Name:    CmdFlag_Output,
				Aliases: []string{"o"},
				Usage:   "output path",
				Value:   ".",
			},
			&cli.StringFlag{
				Name:  CmdFlag_Database,
				Usage: "database name to export",
			},
			&cli.StringFlag{
				Name:    CmdFlag_Tables,
				Aliases: []string{"t"},
				Usage:   "database tables to export",
			},
			&cli.StringFlag{
				Name:    CmdFlag_Tags,
				Aliases: []string{"T"},
				Usage:   "export tags for golang",
			},
			&cli.StringFlag{
				Name:    CmdFlag_Prefix,
				Aliases: []string{"p"},
				Usage:   "filename prefix",
			},
			&cli.StringFlag{
				Name:    CmdFlag_Suffix,
				Aliases: []string{"s"},
				Usage:   "filename suffix",
			},
			&cli.StringFlag{
				Name:    CmdFlag_Package,
				Aliases: []string{"P"},
				Usage:   "package name",
			},
			&cli.StringFlag{
				Name:  CmdFlag_Without,
				Usage: "exclude columns split by colon",
			},
			&cli.StringFlag{
				Name:    CmdFlag_ReadOnly,
				Aliases: []string{"R"},
				Usage:   "readonly columns split by colon",
			},
			&cli.BoolFlag{
				Name:  CmdFlag_Protobuf,
				Usage: "export protobuf file",
			},
			&cli.StringFlag{
				Name:    CmdFlag_SpecType,
				Aliases: []string{"S"},
				Usage:   "specify column as customized types, e.g 'user.detail=UserDetail, user.data=UserData'",
			},
			&cli.BoolFlag{
				Name:    CmdFlag_EnableDecimal,
				Aliases: []string{"D"},
				Usage:   "decimal as sqlca.Decimal type",
			},
			&cli.StringFlag{
				Name:    CmdFlag_GogoOptions,
				Aliases: []string{"gogo"},
				Usage:   "gogo proto options",
			},
			&cli.BoolFlag{
				Name:    CmdFlag_Merge,
				Aliases: []string{"M"},
				Usage:   "export to one file",
			},
			&cli.StringFlag{
				Name:  CmdFlag_DAO,
				Usage: "generate data access object file",
			},
			&cli.StringFlag{
				Name:    CmdFlag_ImportModels,
				Aliases: []string{"im"},
				Usage:   "project name",
			},
			&cli.BoolFlag{
				Name:    CmdFlag_OmitEmpty,
				Aliases: []string{"E"},
				Usage:   "json omitempty",
			},
			&cli.StringFlag{
				Name:    CmdFlag_JsonProperties,
				Aliases: []string{"jp"},
				Usage:   "customized properties for json tag",
			},
			&cli.StringFlag{
				Name:    CmdFlag_TinyintAsBool,
				Aliases: []string{"B"},
				Usage:   "convert tinyint columns redeclare as bool type",
			},
			&cli.StringFlag{
				Name:  CmdFlag_SSH,
				Usage: "ssh tunnel e.g ssh://root:123456@192.168.1.23:22",
			},
			&cli.BoolFlag{
				Name:  CmdFlag_V2,
				Usage: "sqlca v2 package imports",
			},
			&cli.StringFlag{
				Name:    CmdFlag_Export,
				Aliases: []string{"ddl"},
				Usage:   "export database DDL to file",
			},
			&cli.BoolFlag{
				Name:    CmdFlag_Debug,
				Aliases: []string{"d"},
				Usage:   "open debug mode",
			},
			&cli.StringFlag{
				Name:    CmdFlag_ProtoOptions,
				Aliases: []string{"po"},
				Usage:   "set protobuf options, multiple options seperated by ';'",
			},
			&cli.StringFlag{
				Name:    CmdFlag_FieldStyle,
				Aliases: []string{"style"},
				Usage:   "protobuf message field camel style (small or big)",
			},
			&cli.StringFlag{
				Name:    CmdFlag_BaseModel,
				Aliases: []string{"bm"},
				Usage:   "specify base model. e.g github.com/civet148/db2go/types.BaseModel=create_time,update_time,is_deleted",
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
	var cmd = schema.NewCmdFlags()
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
	cmd.OneFile = ctx.Bool(CmdFlag_Merge)
	cmd.JsonProperties = ctx.String(CmdFlag_JsonProperties)
	cmd.ParseSpecTypes(ctx.String(CmdFlag_SpecType))
	cmd.ParseBaseModel(ctx.String(CmdFlag_BaseModel))
	cmd.JsonStyle = ctx.String(CmdFlag_JsonStyle)
	cmd.ExportDDL = ctx.String(CmdFlag_Export)
	cmd.FieldStyle = schema.FieldStyleFromString(ctx.String(CmdFlag_FieldStyle))

	if ctx.Bool(CmdFlag_V2) {
		cmd.SqlcaPkg = schema.SQLCA_V2_PKG
		cmd.ImportVer = schema.IMPORT_SQLCA_V2
	} else {
		cmd.SqlcaPkg = schema.SQLCA_V3_PKG
		cmd.ImportVer = schema.IMPORT_SQLCA_V3
	}
	if cmd.DAO != "" && cmd.ImportModels == "" {
		return log.Errorf("models path required eg. github.com/xxx/your-repo/models")
	}
	if cmd.SSH != "" {
		if !strings.Contains(cmd.SSH, SshScheme) {
			cmd.SSH = SshScheme + cmd.SSH
		}
	}
	if cmd.Protobuf {
		var strProtoOptions string
		strProtoOptions = ctx.String(CmdFlag_ProtoOptions)
		if strProtoOptions != "" {
			opts := strings.Split(strProtoOptions, ";")
			for _, opt := range opts {
				ss := strings.Split(opt, "=")
				if len(ss) != 2 {
					return log.Errorf("invalid protobuf option %s", opt)
				}
				cmd.ProtoOptions[ss[0]] = ss[1]
			}
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

	if v := ctx.String(CmdFlag_Tables); v != "" {
		cmd.Tables = schema.TrimSpaceSlice(schema.Split(v))
	}

	if v := ctx.String(CmdFlag_Without); v != "" {
		cmd.Without = schema.TrimSpaceSlice(schema.Split(v))
	}

	if v := ctx.String(CmdFlag_TinyintAsBool); v != "" {
		cmd.TinyintAsBool = schema.TrimSpaceSlice(schema.Split(v))
	}

	if cmd.Protobuf {
		gogoOpt := ctx.String(CmdFlag_GogoOptions)
		if gogoOpt != "" {
			cmd.GogoOptions = schema.TrimSpaceSlice(schema.Split(gogoOpt))
		}
	}

	if v := ctx.String(CmdFlag_Tags); v != "" {
		cmd.ExtraTags = schema.TrimSpaceSlice(schema.Split(v))
	}
	if v := ctx.String(CmdFlag_ReadOnly); v != "" {
		cmd.ReadOnly = schema.TrimSpaceSlice(schema.Split(v))
	}

	cmd.Scheme = ui.Scheme
	cmd.Host = ui.Host
	cmd.User = ui.User
	cmd.Password = ui.Password

	log.Json("command options", cmd)

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

func export(cmd *schema.CmdFlags, e *sqlca.Engine) (err error) {
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
