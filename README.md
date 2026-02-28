# db2go is a command to export database table structure to go or proto file 

## Usage
```shell
VERSION:
   v3.5.14 20260228 11:38:23 commit 426c24c

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

rem 数据源连接串
set DSN_URL="mysql://root:12345678@127.0.0.1:3306/test?charset=utf8mb4&collation=utf8mb4_unicode_ci&parseTime=true&loc=Local"
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
        rem 安装失败，Linux/Mac请安装gcc工具链，Windows系统可以通过链接直接下载二进制(https://github.com/civet148/release/tree/master/db2go)
        echo error: db2go install failed, Linux/Mac please install gcc tool-chain, Windows download from https://github.com/civet148/release/tree/master/db2go
    )
)


rem 判断db2go是否安装成功
If "%errorlevel%" == "0" (
db2go --debug --url "%DSN_URL%" --out "%OUT_DIR%" --table "%TABLE_NAME%" --json-properties "%JSON_PROPERTIES%" --enable-decimal  --spec-type "%SPEC_TYPES%" ^
--suffix "%SUFFIX_NAME%" --package "%PACK_NAME%" --readonly "%READ_ONLY%" --without "%WITH_OUT%" --tinyint-as-bool "%TINYINT_TO_BOOL%" ^
--tag "%TAGS%" --dao "%DAO_OUT%" --import-models "%IMPORT_MODELS%" --export "%DEPLOY_SQL%"

echo generate go file ok, formatting...
gofmt -w %OUT_DIR%/%PACK_NAME%
)
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
# 指定表名(不指定则整个数据库全部导出)
TABLE_NAME=""
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
# 基础模型声明(指定基础模型类型和字段)
BASE_MODEL="github.com/civet148/db2go/types.BaseModel=create_time,update_time"
# 指定生成数据库建表SQL输出文件路径
DEPLOY_SQL="deploy/test.sql"

# 检查 db2go 是否已安装
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

db2go --debug --url "$DSN_URL" --out "$OUT_DIR" --table "$TABLE_NAME" --json-properties "$JSON_PROPERTIES" --enable-decimal  --spec-type "$SPEC_TYPES" \
--package "$PACK_NAME" --readonly "$READ_ONLY" --without "$WITH_OUT" --dao dao --tinyint-as-bool "$TINYINT_TO_BOOL" \
--tag "$TAGS" --import-models "$IMPORT_MODELS" --export "$DEPLOY_SQL" --base-model "$BASE_MODEL"

echo "generate go file ok, formatting..."
gofmt -w $OUT_DIR/$PACK_NAME

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

