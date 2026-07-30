package schema

import (
	"strings"

	"github.com/civet148/sqlca/v2"
)

type Option func(*CmdFlags)

func WithURL(url string) Option {
	return func(c *CmdFlags) {
		c.ConnUrl = url
	}
}

func WithOutput(out string) Option {
	return func(c *CmdFlags) {
		c.OutDir = out
	}
}

func WithDatabase(db string) Option {
	return func(c *CmdFlags) {
		c.Database = db
	}
}

func WithTables(tables []string) Option {
	return func(c *CmdFlags) {
		c.Tables = tables
	}
}

func WithExcludeTables(exclude []string) Option {
	return func(c *CmdFlags) {
		c.ExcludeTables = exclude
	}
}

func WithTableFilter(table string) Option {
	return func(c *CmdFlags) {
		if table == "" {
			return
		}
		tables := TrimSpaceSlice(Split(table))
		for _, t := range tables {
			if t[0] == '-' {
				c.ExcludeTables = append(c.ExcludeTables, t[1:])
			} else {
				c.Tables = append(c.Tables, t)
			}
		}
	}
}

func WithTags(tags []string) Option {
	return func(c *CmdFlags) {
		c.ExtraTags = tags
		var hasGorm bool
		for _, tag := range tags {
			if tag == "gorm" {
				hasGorm = true
			}
		}
		if !hasGorm {
			c.ExtraTags = append(c.ExtraTags, "gorm")
		}
	}
}

func WithPrefix(prefix string) Option {
	return func(c *CmdFlags) {
		c.Prefix = prefix
	}
}

func WithSuffix(suffix string) Option {
	return func(c *CmdFlags) {
		c.Suffix = suffix
	}
}

func WithPackageName(pkg string) Option {
	return func(c *CmdFlags) {
		c.PackageName = pkg
	}
}

func WithWithout(without []string) Option {
	return func(c *CmdFlags) {
		c.Without = without
	}
}

func WithReadOnly(readonly []string) Option {
	return func(c *CmdFlags) {
		c.ReadOnly = readonly
	}
}

func WithSpecType(specType string) Option {
	return func(c *CmdFlags) {
		c.ParseSpecTypes(specType)
	}
}

func WithEnableDecimal(enable bool) Option {
	return func(c *CmdFlags) {
		c.EnableDecimal = enable
	}
}

func WithSSH(ssh string) Option {
	return func(c *CmdFlags) {
		if ssh != "" {
			if !strings.Contains(ssh, "ssh://") {
				c.SSH = "ssh://" + ssh
			} else {
				c.SSH = ssh
			}
		}
	}
}

func WithV2(v2 bool) Option {
	return func(c *CmdFlags) {
		if v2 {
			c.SqlcaPkg = SQLCA_V2_PKG
			c.ImportVer = IMPORT_SQLCA_V2
		} else {
			c.SqlcaPkg = SQLCA_V3_PKG
			c.ImportVer = IMPORT_SQLCA_V3
		}
	}
}

func WithDebug(debug bool) Option {
	return func(c *CmdFlags) {
		c.Debug = debug
	}
}

func WithFieldStyle(style string) Option {
	return func(c *CmdFlags) {
		c.FieldStyle = FieldStyleFromString(style)
	}
}

func WithExportDDL(ddl string) Option {
	return func(c *CmdFlags) {
		c.ExportDDL = ddl
	}
}

func WithJsonStyle(style string) Option {
	return func(c *CmdFlags) {
		c.JsonStyle = style
	}
}

func WithProtoOptions(opts map[string]string) Option {
	return func(c *CmdFlags) {
		c.ProtoOptions = opts
	}
}

func WithGogoOptions(opts []string) Option {
	return func(c *CmdFlags) {
		c.GogoOptions = opts
	}
}

func WithOneFile(merge bool) Option {
	return func(c *CmdFlags) {
		c.OneFile = merge
	}
}

func WithEngine(e *sqlca.Engine) Option {
	return func(c *CmdFlags) {
		c.Engine = e
	}
}

func WithScheme(scheme string) Option {
	return func(c *CmdFlags) {
		c.Scheme = scheme
	}
}

func WithProtobuf(proto bool) Option {
	return func(c *CmdFlags) {
		c.Protobuf = proto
	}
}

func WithTagTypes(tagTypes []*CommTagType) Option {
	return func(c *CmdFlags) {
		c.TagTypes = tagTypes
	}
}
