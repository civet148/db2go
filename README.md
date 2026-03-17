# db2go is a command to export database table structure to go or proto file 

## Usage
```shell
NAME:
   db2go - db2go [options] --url <DSN>

USAGE:
   db2go [global options] command [command options] [arguments...]

VERSION:
   v3.6.0 20260317 17:47:51 commit ffd2e42

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --url value, -u value                data source name of database
   --out value, -o value                output path (default: ".")
   --db value                           database name to export
   --table value, -t value              database tables to export
   --tag value, -T value                export tags for golang
   --prefix value, -p value             filename prefix
   --suffix value, -s value             filename suffix
   --package value, -P value            package name
   --without value                      exclude columns split by colon
   --readonly value, -R value           readonly columns split by colon
   --proto                              export protobuf file (default: false)
   --spec-type value, -S value          specify column as customized types, e.g 'user.detail=UserDetail, user.data=UserData'
   --enable-decimal, -D                 decimal as sqlca.Decimal type (default: false)
   --gogo-options value, --gogo value   gogo proto options
   --merge, -M                          export to one file (default: false)
   --dao value                          generate data access object file
   --import-models value, --im value    project name
   --omitempty, -E                      json omitempty (default: false)
   --json-properties value, --jp value  customized properties for json tag
   --tinyint-as-bool value, -B value    convert tinyint columns redeclare as bool type
   --ssh value                          ssh tunnel e.g ssh://root:123456@192.168.1.23:22
   --v2                                 sqlca v2 package imports (default: false)
   --export value, --ddl value          export database DDL to file
   --debug, -d                          open debug mode (default: false)
   --proto-options value, --po value    set protobuf options, multiple options seperated by ';'
   --field-style value, --style value   protobuf message field camel style (small or big)
   --base-model value, --bm value       specify base model. e.g types.BaseModel=created_at,updated_at
   --preload-model value, --pm value    specify gorm preload model. e.g users.Roles=[]*Role(many2many:user_roles), users.Profile=UserProfile(foreignKey:UserId;)t
   --help, -h                           show help (default: false)
   --version, -v                        print the version (default: false)

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
rem tinyint转换成bool
set TINYINT_TO_BOOL="deleted,is_admin,disable"
rem 数据库连接源DSN
set DSN_URL="mysql://root:123456@127.0.0.1:3306/test?charset=utf8"
rem JSON属性
set JSON_PROPERTIES="omitempty"
rem 指定具体表对应字段类型(不指定表则全局生效)
set SPEC_TYPES="users.extra_data=struct%%, users.is_deleted=bool"
rem 导入models路径(仅生成DAO文件使用)
set IMPORT_MODELS="github.com/civet148/db2go/models"
rem 基础模型声明
set BASE_MODEL="BaseModel=created_at,updated_at"
rem 预加载模型声明
set PRELOAD_MODEL="users.Roles=[]*Role(many2many:user_roles), users.Profile=UserProfile(foreignKey:UserId;)"
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


db2go --url "%DSN_URL%" --out "%OUT_DIR%" --table "%TABLE_NAME%" --json-properties "%JSON_PROPERTIES%" --enable-decimal  --spec-type "%SPEC_TYPES%" \
 --package "%PACK_NAME%" --readonly "%READ_ONLY%" --without "%WITH_OUT%" --dao dao --tinyint-as-bool "%TINYINT_TO_BOOL%" \
 --tag "%TAGS%" --import-models %IMPORT_MODELS% --base-model "%BASE_MODEL%" --ddl "%DDL_FILE%" --preload-model "%PRELOAD_MODEL%"

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
# TINYINT转换成bool
TINYINT_TO_BOOL="deleted,is_admin,disable"
# 数据库连接源DSN
DSN_URL="mysql://root:123456@127.0.0.1:3306/test?charset=utf8"
# JSON属性
JSON_PROPERTIES="omitempty"
# 指定具体表对应字段类型(不指定表则全局生效)
SPEC_TYPES="users.extra_data=struct{}, users.is_deleted=bool"
# 导入models路径(仅生成DAO文件使用)
IMPORT_MODELS="github.com/civet148/db2go/models"
# 基础模型声明
BASE_MODEL="BaseModel=created_at,updated_at"
# 预加载模型声明
PRELOAD_MODEL="users.Roles=[]*Role(many2many:user_roles), users.Profile=UserProfile(foreignKey:UserId;)"
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

db2go --url "${DSN_URL}" --out "${OUT_DIR}" --table "${TABLE_NAME}" --json-properties "${JSON_PROPERTIES}" --enable-decimal  --spec-type "${SPEC_TYPES}" \
 --package "${PACK_NAME}" --readonly "${READ_ONLY}" --without "${WITH_OUT}" --dao dao --tinyint-as-bool "${TINYINT_TO_BOOL}" \
 --tag "${TAGS}" --import-models ${IMPORT_MODELS} --base-model "${BASE_MODEL}" --ddl "${DDL_FILE}" --preload-model "${PRELOAD_MODEL}"

echo "generate go file ok, formatting..."
gofmt -w ${OUT_DIR}/${PACK_NAME}

```

- data object

```go
package models

import github_com_civet148_db2go_types "github.com/civet148/db2go/types"
import "github.com/civet148/sqlca/v3"

const TableNameInventoryData = "inventory_data" //产品库存数据表

const (
	INVENTORY_DATA_COLUMN_ID            = "id"
	INVENTORY_DATA_COLUMN_CREATE_ID     = "create_id"
	INVENTORY_DATA_COLUMN_CREATE_NAME   = "create_name"
	INVENTORY_DATA_COLUMN_CREATE_TIME   = "create_time"
	INVENTORY_DATA_COLUMN_UPDATE_ID     = "update_id"
	INVENTORY_DATA_COLUMN_UPDATE_NAME   = "update_name"
	INVENTORY_DATA_COLUMN_UPDATE_TIME   = "update_time"
	INVENTORY_DATA_COLUMN_IS_FROZEN     = "is_frozen"
	INVENTORY_DATA_COLUMN_NAME          = "name"
	INVENTORY_DATA_COLUMN_SERIAL_NO     = "serial_no"
	INVENTORY_DATA_COLUMN_QUANTITY      = "quantity"
	INVENTORY_DATA_COLUMN_PRICE         = "price"
	INVENTORY_DATA_COLUMN_PRODUCT_EXTRA = "product_extra"
	INVENTORY_DATA_COLUMN_LOCATION      = "location"
)

type InventoryData struct {
	github_com_civet148_db2go_types.BaseModel
	Id           uint64        `json:"id,omitempty" db:"id" gorm:"column:id;primaryKey;autoIncrement;"`                                  //产品ID
	CreateId     uint64        `json:"create_id,omitempty" db:"create_id" gorm:"column:create_id;type:bigint unsigned;default:0;"`       //创建人ID
	CreateName   string        `json:"create_name,omitempty" db:"create_name" gorm:"column:create_name;type:varchar(64);"`               //创建人姓名
	UpdateId     uint64        `json:"update_id,omitempty" db:"update_id" gorm:"column:update_id;type:bigint unsigned;default:0;"`       //更新人ID
	UpdateName   string        `json:"update_name,omitempty" db:"update_name" gorm:"column:update_name;type:varchar(64);"`               //更新人姓名
	IsFrozen     int8          `json:"is_frozen,omitempty" db:"is_frozen" gorm:"column:is_frozen;type:tinyint(1);default:0;"`            //冻结状态(0: 未冻结 1: 已冻结)
	Name         string        `json:"name,omitempty" db:"name" gorm:"column:name;type:varchar(255);"`                                   //产品名称
	SerialNo     string        `json:"serial_no,omitempty" db:"serial_no" gorm:"column:serial_no;type:varchar(64);"`                     //产品编号
	Quantity     sqlca.Decimal `json:"quantity,omitempty" db:"quantity" gorm:"column:quantity;type:decimal(16,3);default:0.000;"`        //产品库存
	Price        sqlca.Decimal `json:"price,omitempty" db:"price" gorm:"column:price;type:decimal(16,2);default:0.00;"`                  //产品均价
	ProductExtra string        `json:"product_extra,omitempty" db:"product_extra" gorm:"column:product_extra;type:text;" sqlca:"isnull"` //产品附带数据(JSON文本)
	Location     sqlca.Point   `json:"location,omitempty" db:"location" gorm:"column:location;type:point;" sqlca:"isnull"`               //地理位置
}

func (do *InventoryData) GetId() uint64               { return do.Id }
func (do *InventoryData) SetId(v uint64)              { do.Id = v }
func (do *InventoryData) GetCreateId() uint64         { return do.CreateId }
func (do *InventoryData) SetCreateId(v uint64)        { do.CreateId = v }
func (do *InventoryData) GetCreateName() string       { return do.CreateName }
func (do *InventoryData) SetCreateName(v string)      { do.CreateName = v }
func (do *InventoryData) GetCreateTime() string       { return do.CreateTime }
func (do *InventoryData) SetCreateTime(v string)      { do.CreateTime = v }
func (do *InventoryData) GetUpdateId() uint64         { return do.UpdateId }
func (do *InventoryData) SetUpdateId(v uint64)        { do.UpdateId = v }
func (do *InventoryData) GetUpdateName() string       { return do.UpdateName }
func (do *InventoryData) SetUpdateName(v string)      { do.UpdateName = v }
func (do *InventoryData) GetUpdateTime() string       { return do.UpdateTime }
func (do *InventoryData) SetUpdateTime(v string)      { do.UpdateTime = v }
func (do *InventoryData) GetIsFrozen() int8           { return do.IsFrozen }
func (do *InventoryData) SetIsFrozen(v int8)          { do.IsFrozen = v }
func (do *InventoryData) GetName() string             { return do.Name }
func (do *InventoryData) SetName(v string)            { do.Name = v }
func (do *InventoryData) GetSerialNo() string         { return do.SerialNo }
func (do *InventoryData) SetSerialNo(v string)        { do.SerialNo = v }
func (do *InventoryData) GetQuantity() sqlca.Decimal  { return do.Quantity }
func (do *InventoryData) SetQuantity(v sqlca.Decimal) { do.Quantity = v }
func (do *InventoryData) GetPrice() sqlca.Decimal     { return do.Price }
func (do *InventoryData) SetPrice(v sqlca.Decimal)    { do.Price = v }
func (do *InventoryData) GetProductExtra() string     { return do.ProductExtra }
func (do *InventoryData) SetProductExtra(v string)    { do.ProductExtra = v }
func (do *InventoryData) GetLocation() sqlca.Point    { return do.Location }
func (do *InventoryData) SetLocation(v sqlca.Point)   { do.Location = v }

////////////////////// ----- 自定义代码请写在下面 ----- //////////////////////

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
    db2go --url %DSN_URL% --proto --out %OUT_DIR% --table %TABLE_NAME% --package %PACK_NAME%  --without %WITH_OUT% --proto-options %PROTO_OPTION% --merge
    echo generate protobuf files ok
    gofmt -w %OUT_DIR%/%PACK_NAME%
)

PAUSE

```

