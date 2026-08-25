# db2go is a command to export database table structure to go or proto file 

## Usage
```shell
NAME:
   db2go - db2go [command] [options]

USAGE:
   db2go [global options] command [command options] [arguments...]

VERSION:
   v3.9.1 20260731 10:29:06 commit 67292be

COMMANDS:
   go       Generate Go model files from database tables
   proto    Generate protobuf files from database tables
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help (default: false)
   --version, -v  print the version (default: false)
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

rem 输出文件根目录
set OUT_DIR=.
rem 数据模型文件包名
set PACK_NAME="models"
rem 只读字段(不更新)
set READ_ONLY="created_at, updated_at"
rem 指定或排除表名(不指定则整个数据库全部导出, 排除表名在表名前面加-)
set TABLE_NAME="-user_roles"
rem 忽略字段名(逗号分隔)
set WITH_OUT=""
rem 添加标签
set TAGS="gorm"
rem 数据库连接源DSN
set DSN_URL="mysql://root:123456@127.0.0.1:3306/test?charset=utf8"
rem 指定具体表对应字段类型(不指定表则全局生效)
set SPEC_TYPES="users.extra_data=struct%%, users.is_deleted=bool"
rem 数据库DDL文件
set DDL_FILE="deploy/test.sql"


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
        rem 安装失败，Linux/Mac请安装gcc工具链，Windows系统可以通过链接直接下载二进制(https://github.com/civet148/release/tree/master/db2go)
        echo error: db2go install failed, Linux/Mac please install gcc tool-chain, Windows download from https://github.com/civet148/release/tree/master/db2go
    )
)


db2go go --url "%DSN_URL%" --out "%OUT_DIR%" --table "%TABLE_NAME%" --enable-decimal  --spec-type "%SPEC_TYPES%" ^
 --package "%PACK_NAME%" --readonly "%READ_ONLY%" --without "%WITH_OUT%" --tag "%TAGS%" --ddl "%DDL_FILE%" 

echo "generate go file ok, formatting..."
gofmt -w %OUT_DIR%/%PACK_NAME%
pause

```
- Linux/Unix shell脚本

```shell
#!/bin/sh

# 输出文件根目录
OUT_DIR=.
# 数据模型文件包名
PACK_NAME="models"
# 只读字段(不更新)
READ_ONLY="created_at, updated_at"
# 指定或排除表名(不指定则整个数据库全部导出, 排除表名在表名前面加-)
TABLE_NAME="-user_roles"
# 忽略字段名(逗号分隔)
WITH_OUT=""
# 添加标签
TAGS="gorm"
# 数据库连接源DSN
DSN_URL="mysql://root:123456@127.0.0.1:3306/test?charset=utf8"
# 指定具体表对应字段类型(不指定表则全局生效)
SPEC_TYPES="users.extra_data=struct{}, users.is_deleted=bool"
# 数据库DDL文件
DDL_FILE="deploy/test.sql"

## 检查 db2go 是否已安装
if ! which db2go >/dev/null 2>&1; then
    # 安装最新版 db2go
    go install github.com/civet148/db2go@latest

    # 检查是否安装成功
    if which db2go >/dev/null 2>&1; then
        echo "✅ db2go install success, $(which db2go)"
    else
        echo "❌ db2go install failed, please check go env and gcc tool-chain"
        exit 1
    fi
fi

db2go go --url "${DSN_URL}" --out "${OUT_DIR}" --table "${TABLE_NAME}" --enable-decimal  --spec-type "${SPEC_TYPES}" \
 --package "${PACK_NAME}" --readonly "${READ_ONLY}" --without "${WITH_OUT}" --tag "${TAGS}" --ddl "${DDL_FILE}" 

echo "generate go file ok, formatting..."
gofmt -w ${OUT_DIR}/${PACK_NAME}

```

- data object

```go
package models

import (
    "time"
)


// 手动定义的用户额外数据结构，对应数据库表 users 的 extra_data 字段
// db2go自动合并数据库导出部分代码和手动写的代码
type UserExtraData struct {
	HomeAddress string `json:"home_address"`
	PostCode    string `json:"post_code"`
}

const TableNameUsers = "users"

const (
    UsersColumn_Id        = "id"
    UsersColumn_CreatedAt = "created_at"
    UsersColumn_UpdatedAt = "updated_at"
    UsersColumn_UserName  = "user_name"
    UsersColumn_Email     = "email"
    UsersColumn_ExtraData = "extra_data"
)

type User struct {
    Id          uint64        `json:"id" db:"id" gorm:"column:id;primaryKey;autoIncrement;"`
    CreatedAt   time.Time     `json:"created_at" db:"created_at" gorm:"column:created_at;type:timestamp;autoCreateTime;index:idx_users_created_at;default:CURRENT_TIMESTAMP;" sqlca:"readonly"`
    UpdatedAt   time.Time     `json:"updated_at" db:"updated_at" gorm:"column:updated_at;type:timestamp;autoUpdateTime;index:idx_users_updated_at;default:CURRENT_TIMESTAMP;" sqlca:"readonly"`
    UserName    string        `json:"user_name" db:"user_name" gorm:"column:user_name;type:varchar(32);uniqueIndex:idx_users_user_name;default:null;" sqlca:"isnull"`
    Email       string        `json:"email" db:"email" gorm:"column:email;type:varchar(64);uniqueIndex:idx_users_email;default:null;" sqlca:"isnull"`
    ExtraData   UserExtraData `json:"extra_data" db:"extra_data" gorm:"column:extra_data;type:json;default:null;" sqlca:"isnull"`
	Roles       []*Role       `json:"roles,omitempty" db:"-" gorm:"many2many:user_roles;-:migration;"` // 用户角色列表(手工添加，自动合并)
	Profile     UserProfile   `json:"profile,omitempty" db:"-" gorm:"foreignKey:UserId;-:migration;"`  // 用户资料明细(手工添加，自动合并)
}

func (do User) TableName() string {
    return TableNameUsers
}

func (do User) GetId() uint64 { return do.Id }
func (do User) GetCreatedAt() time.Time { return do.CreatedAt }
func (do User) GetUpdatedAt() time.Time { return do.UpdatedAt }
func (do User) GetUserName() string { return do.UserName }
func (do User) GetEmail() string { return do.Email }
func (do User) GetExtraData() UserExtraData { return do.ExtraData }

func (do *User) SetId(v uint64) { do.Id = v }
func (do *User) SetCreatedAt(v time.Time) { do.CreatedAt = v }
func (do *User) SetUpdatedAt(v time.Time) { do.UpdatedAt = v }
func (do *User) SetUserName(v string) { do.UserName = v }
func (do *User) SetEmail(v string) { do.Email = v }
func (do *User) SetExtraData(v UserExtraData) { do.ExtraData = v }

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
        rem 安装失败，Linux/Mac请安装gcc工具链，Windows系统可以通过链接直接下载二进制(https://github.com/civet148/release/tree/master/db2go)
        echo error: db2go install failed, Linux/Mac please install gcc tool-chain, Windows download from https://github.com/civet148/release/tree/master/db2go
    )
)

rem 判断db2go是否安装成功
IF "%errorlevel%" == "0" (
    db2go --url %DSN_URL% --proto --out %OUT_DIR% --table %TABLE_NAME% --package %PACK_NAME%  --without %WITH_OUT% --proto-options %PROTO_OPTION%
    echo generate protobuf files ok
    gofmt -w %OUT_DIR%/%PACK_NAME%
)

PAUSE

```

