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
# 基础模型声明
BASE_MODEL="github.com/civet148/db2go/types.BaseModel=create_time,update_time"
# 数据库DDL文件
DDL_FILE="deploy/test.sql"

rm -rf ./models ./dao

./db2go --url "$DSN_URL" --out "$OUT_DIR" --table "$TABLE_NAME" --json-properties "$JSON_PROPERTIES" --enable-decimal  --spec-type "$SPEC_TYPES" \
 --package "$PACK_NAME" --readonly "$READ_ONLY" --without "$WITH_OUT" --dao dao --tinyint-as-bool "$TINYINT_TO_BOOL" \
 --tag "$TAGS" --import-models $IMPORT_MODELS --base-model "$BASE_MODEL" --ddl "$DDL_FILE"

echo "generate go file ok, formatting..."
gofmt -w $OUT_DIR/$PACK_NAME
