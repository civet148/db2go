# db2go is a command to export database table structure to go or proto file 

## Usage
```shell
$ db2go -h
NAME:
   db2go - db2go [options] --url <DSN>

USAGE:
   db2go [global options] command [command options] [arguments...]

VERSION:
   v2.13.0 2024-07-31 commit <N/A>

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --url value              data source name of database
   --out value              output path (default: ".")
   --db value               database name to export
   --table value            database tables to export
   --tag value              export tags for golang
   --prefix value           filename prefix
   --suffix value           filename suffix
   --package value          package name
   --without value          exclude columns split by colon
   --readonly value         readonly columns split by colon
   --proto                  export protobuf file (default: false)
   --spec-type value        specify column as customized types, e.g 'user.detail=UserDetail, user.data=UserData'
   --enable-decimal         decimal as sqlca.Decimal type (default: false)
   --gogo-options value     gogo proto options
   --merge                  export to one file (default: false)
   --dao value              generate data access object file
   --import-models value    project name
   --omitempty              json omitempty (default: false)
   --json-properties value  customized properties for json tag
   --tinyint-as-bool value  convert tinyint columns redeclare as bool type
   --ssh value              ssh tunnel e.g ssh://root:123456@192.168.1.23:22
   --v1                     v1 package imports (default: false)
   --export value           export database DDL to file
   --debug                  open debug mode (default: false)
   --json-style value       json style: smallcamel or bigcamel (default: "default")
   --common-tags value      set common tag values, e.g gorm
   --proto-options value    set protobuf options, multiple options seperated by ';'
   --help, -h               show help (default: false)
   --version, -v            print the version (default: false)

```

## 1. 编译安装

- Ubuntu 20.04 or later

```shell
$ sudo apt update && sudo apt install -y make gcc 
$ go env -w CGO_ENABLED=1
$ make
```

## 2. 数据库表导出到go文件

* Windows batch 脚本

```batch
@echo off

rem 数据源连接串
set DSN_URL="mysql://root:123456@127.0.0.1:3306/test?charset=utf8mb4&collation=utf8mb4_unicode_ci&parseTime=true&loc=Local"
rem 数据模型(models)和数据库操作对象(dao)文件输出基础目录
set OUT_DIR=.
rem 数据模型包名(数据模型文件目录名)
set PACK_NAME="models"
rem 文件名后缀(一般不需要加)
set SUFFIX_NAME=""
rem 只读字段（例如create_at/update_at这种通过设置数据库自动生成/更新的时间字段）,使用sqlca包Insert/Update/Upsert方法时不会对这些字段进行插入和更新操作）
set READ_ONLY=""
rem 指定数据库表名（留空表示导出全部表结构）
set TABLE_NAME=""
rem 忽略哪些数据库表字段
set WITH_OUT=""
rem 增加自定义标签名（例如: bson）
set TAGS=""
rem 指定所有表结构中属于tinyint类型的某些字段为布尔类型
set TINYINT_TO_BOOL=""
rem 附加的json标签属性（例如: omitempty)
set JSON_PROPERTIES=""
rem 指定某表的某字段为指定类型,多个表字段以英文逗号分隔（例如：user.is_deleted=bool表示指定user表is_deleted字段为bool类型; 如果不指定表名则所有表的is_deleted字段均为bool类型；支持第三方包类型，例如：user.weight=github.com/shopspring/decimal.Decimal）
set SPEC_TYPES="users.extra_data=struct{}"
rem 指定其他orm的标签和值(以空格分隔)
set COMMON_TAGS="id=gorm:primarykey created_at=gorm:\"autoCreateTime;type:timestamp\" updated_at=gorm:\"autoUpdateTime;type:timestamp\""
rem 数据库操作对象生成目录名
set DAO_OUT=dao
rem 数据库操作对象导入数据库表模型数据路径(指定--dao选项时必填)
set IMPORT_MODELS="test/sqler/models"
rem 导出全部建表SQL到指定文件
set DEPLOY_SQL="deploy/test.sql"

rem 判断本地系统是否已安装db2go工具，没有则进行安装
where db2go.exe

IF "%errorlevel%" == "0" (
    echo db2go already installed.
) ELSE (
    echo db2go not found in system %%PATH%%, installing...
    go install github.com/civet148/db2go@latest
    If "%errorlevel%" == "0" (
        echo db2go install succeeded
    ) ELSE (
        rem 安装失败，Linux/Mac请安装gcc工具链，Windows系统可以通过链接直接下载二进制(https://github.com/civet148/release/tree/master/db2go/v2)
        echo error: db2go install failed, Linux/Mac please install gcc tool-chain, Windows download from https://github.com/civet148/release/tree/master/db2go/v2
    )
)


rem 判断db2go是否安装成功
If "%errorlevel%" == "0" (
db2go --url %DSN_URL% --out %OUT_DIR% --table %TABLE_NAME% --json-properties %JSON_PROPERTIES% --enable-decimal  --spec-type %SPEC_TYPES% ^
--suffix %SUFFIX_NAME% --package %PACK_NAME% --readonly %READ_ONLY% --without %WITH_OUT% --tinyint-as-bool %TINYINT_TO_BOOL% ^
--tag %TAGS% --common-tags %COMMON_TAGS% --dao %DAO_OUT% --import-models %IMPORT_MODELS% ^
rem --export %DEPLOY_SQL%

echo generate go file ok, formatting...
gofmt -w %OUT_DIR%/%PACK_NAME%
)
pause

```
- Linux/Unix shell脚本

```shell
#!/bin/sh

# 数据对象文件输出目录
OUT_DIR=.
# 数据对象文件golang包名
PACK_NAME="models"
# 文件后缀名(一般不需要加)
SUFFIX_NAME=""
# 指定只读字段(当created_at和updated_at字段由数据库自动生成时使用)
READ_ONLY="created_at, updated_at"
# 数据库表名(可为空)
TABLE_NAME="users, classes"
# 排除某些字段
WITH_OUT=""
# 字段自定义标签
TAGS="bson"
# 指定字段自动生成为bool类型
TINYINT_TO_BOOL="deleted,is_admin,disable"
# 数据源连接URL
DSN_URL="mysql://root:123456@127.0.0.1:3306/test?charset=utf8"
# 增加json标签属性
JSON_PROPERTIES="omitempty"
# 指定字段为自定义类型(如果是结构体则需要自行创建一个文件进行声明)
SPEC_TYPES="users.extra_data=struct{}"
# 自动生成数据库操作对象文件时需指定数据对象文件导入路径
IMPORT_MODELS="github.com/civet148/db2go/models"

if [ $? -eq 0 ]; then
db2go --debug --url $DSN_URL --out $OUT_DIR --table "$TABLE_NAME" --json-properties $JSON_PROPERTIES --enable-decimal  --spec-type "$SPEC_TYPES" \
--suffix $SUFFIX_NAME --package $PACK_NAME --readonly "$READ_ONLY" --without "$WITH_OUT" --dao dao --tinyint-as-bool "$TINYINT_TO_BOOL" \
--tag "$TAGS" --import-models $IMPORT_MODELS

echo "generate go file ok, formatting..."
gofmt -w $OUT_DIR/$PACK_NAME

else
  echo "error: db2go build failed"
fi
```

- data object

```go
// Code generated by db2go. DO NOT EDIT.
// https://github.com/civet148/sqlca

package models

import "github.com/civet148/sqlca/v2"

const TableNameUsers = "users" //

const (
	USERS_COLUMN_ID         = "id"
	USERS_COLUMN_NAME       = "name"
	USERS_COLUMN_PHONE      = "phone"
	USERS_COLUMN_SEX        = "sex"
	USERS_COLUMN_EMAIL      = "email"
	USERS_COLUMN_DISABLE    = "disable"
	USERS_COLUMN_BALANCE    = "balance"
	USERS_COLUMN_SEX_NAME   = "sex_name"
	USERS_COLUMN_DATA_SIZE  = "data_size"
	USERS_COLUMN_EXTRA_DATA = "extra_data"
	USERS_COLUMN_CREATED_AT = "created_at"
	USERS_COLUMN_UPDATED_AT = "updated_at"
	USERS_COLUMN_DELETED_AT = "deleted_at"
)

type UsersDO struct {
	Id        uint32        `json:"id,omitempty" db:"id" bson:"_id"`                                         //auto inc id
	Name      string        `json:"name,omitempty" db:"name" bson:"name"`                                    //user name
	Phone     string        `json:"phone,omitempty" db:"phone" bson:"phone"`                                 //phone number
	Sex       uint8         `json:"sex,omitempty" db:"sex" bson:"sex"`                                       //user sex
	Email     string        `json:"email,omitempty" db:"email" bson:"email"`                                 //email
	Disable   bool          `json:"disable,omitempty" db:"disable" bson:"disable"`                           //disabled(0=false 1=true)
	Balance   sqlca.Decimal `json:"balance,omitempty" db:"balance" bson:"balance"`                           //balance of decimal
	SexName   string        `json:"sex_name,omitempty" db:"sex_name" bson:"sex_name"`                        //sex name
	DataSize  int64         `json:"data_size,omitempty" db:"data_size" bson:"data_size"`                     //data size
	ExtraData struct{}      `json:"extra_data,omitempty" db:"extra_data" sqlca:"isnull" bson:"extra_data"`   //extra data
	CreatedAt string        `json:"created_at,omitempty" db:"created_at" sqlca:"readonly" bson:"created_at"` //create time
	UpdatedAt string        `json:"updated_at,omitempty" db:"updated_at" sqlca:"readonly" bson:"updated_at"` //update time
	DeletedAt string        `json:"deleted_at,omitempty" db:"deleted_at" sqlca:"isnull" bson:"deleted_at"`   //delete time
}

func (do *UsersDO) GetId() uint32              { return do.Id }
func (do *UsersDO) SetId(v uint32)             { do.Id = v }
func (do *UsersDO) GetName() string            { return do.Name }
func (do *UsersDO) SetName(v string)           { do.Name = v }
func (do *UsersDO) GetPhone() string           { return do.Phone }
func (do *UsersDO) SetPhone(v string)          { do.Phone = v }
func (do *UsersDO) GetSex() uint8              { return do.Sex }
func (do *UsersDO) SetSex(v uint8)             { do.Sex = v }
func (do *UsersDO) GetEmail() string           { return do.Email }
func (do *UsersDO) SetEmail(v string)          { do.Email = v }
func (do *UsersDO) GetDisable() bool           { return do.Disable }
func (do *UsersDO) SetDisable(v bool)          { do.Disable = v }
func (do *UsersDO) GetBalance() sqlca.Decimal  { return do.Balance }
func (do *UsersDO) SetBalance(v sqlca.Decimal) { do.Balance = v }
func (do *UsersDO) GetSexName() string         { return do.SexName }
func (do *UsersDO) SetSexName(v string)        { do.SexName = v }
func (do *UsersDO) GetDataSize() int64         { return do.DataSize }
func (do *UsersDO) SetDataSize(v int64)        { do.DataSize = v }
func (do *UsersDO) GetExtraData() struct{}     { return do.ExtraData }
func (do *UsersDO) SetExtraData(v struct{})    { do.ExtraData = v }
func (do *UsersDO) GetCreatedAt() string       { return do.CreatedAt }
func (do *UsersDO) SetCreatedAt(v string)      { do.CreatedAt = v }
func (do *UsersDO) GetUpdatedAt() string       { return do.UpdatedAt }
func (do *UsersDO) SetUpdatedAt(v string)      { do.UpdatedAt = v }
func (do *UsersDO) GetDeletedAt() string       { return do.DeletedAt }
func (do *UsersDO) SetDeletedAt(v string)      { do.DeletedAt = v }

/*
CREATE TABLE `users` (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT 'auto inc id',
  `name` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT 'user name',
  `phone` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT 'phone number',
  `sex` tinyint unsigned NOT NULL DEFAULT '0' COMMENT 'user sex',
  `email` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT 'email',
  `disable` tinyint(1) NOT NULL DEFAULT '0' COMMENT 'disabled(0=false 1=true)',
  `balance` decimal(10,2) NOT NULL DEFAULT '0.00' COMMENT 'balance of decimal',
  `sex_name` varchar(16) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT 'sex name',
  `data_size` bigint NOT NULL DEFAULT '0' COMMENT 'data size',
  `extra_data` json DEFAULT NULL COMMENT 'extra data',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'create time',
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'update time',
  `deleted_at` datetime DEFAULT NULL COMMENT 'delete time',
  PRIMARY KEY (`id`) USING BTREE,
  KEY `phone` (`phone`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci ROW_FORMAT=DYNAMIC;
*/

```

## 2. 数据库表导出到proto文件

```batch
@echo off

rem 数据源连接串
set DSN_URL="mysql://root:123456@127.0.0.1:3306/test?charset=utf8mb4&collation=utf8mb4_unicode_ci&parseTime=true&loc=Local"
rem 数据模型(models)和数据库操作对象(dao)文件输出基础目录
set OUT_DIR=.
rem 数据模型包名(数据模型文件目录名)
set PACK_NAME="protos"
rem 指定数据库表名（留空表示导出全部表结构）
set TABLE_NAME=""
rem 忽略哪些数据库表字段
set WITH_OUT=""
rem 设置protobuf option
set PROTO_OPTION="go_package=./pb"

rem 判断本地系统是否已安装db2go工具，没有则进行安装
where db2go.exe

IF "%errorlevel%" == "0" (
    echo db2go already installed.
) ELSE (
    echo db2go not found in system %%PATH%%, installing...
    go install github.com/civet148/db2go@latest
    If "%errorlevel%" == "0" (
        echo db2go install succeeded
    ) ELSE (
        rem 安装失败，Linux/Mac请安装gcc工具链，Windows系统可以通过链接直接下载二进制(https://github.com/civet148/release/tree/master/db2go/v2)
        echo error: db2go install failed, Linux/Mac please install gcc tool-chain, Windows download from https://github.com/civet148/release/tree/master/db2go/v2
    )
)

rem 判断db2go是否安装成功
IF "%errorlevel%" == "0" (
    db2go --url %DSN_URL% --proto --out %OUT_DIR% --table %TABLE_NAME% --package %PACK_NAME%  --without %WITH_OUT% --proto-options %PROTO_OPTION% --merge
    echo generate protobuf files ok
    gofmt -w %OUT_DIR%/%PACK_NAME%
)

PAUSE

```

