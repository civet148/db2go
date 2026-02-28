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
# 数据库DDL文件
DDL_FILE="deploy/test.sql"

### 检查 db2go 是否已安装
#if ! which db2go >/dev/null 2>&1; then
#    # 安装最新版 db2go
#    go install github.com/civet148/db2go@latest
#
#    # 检查是否安装成功
#    if which db2go >/dev/null 2>&1; then
#        echo "✅ db2go install success, $(which db2go)"
#    else
#        echo "❌ db2go install failed, please check go env and gcc tool-chain"
#        exit 1
#    fi
#fi

./db2go --url "$DSN_URL" --out "$OUT_DIR" --table "$TABLE_NAME" --json-properties "$JSON_PROPERTIES" --enable-decimal  --spec-type "$SPEC_TYPES" \
 --package "$PACK_NAME" --readonly "$READ_ONLY" --without "$WITH_OUT" --dao dao --tinyint-as-bool "$TINYINT_TO_BOOL" \
 --tag "$TAGS" --import-models $IMPORT_MODELS --base-model "$BASE_MODEL" --ddl "$DDL_FILE"

echo "generate go file ok, formatting..."
gofmt -w $OUT_DIR/$PACK_NAME
