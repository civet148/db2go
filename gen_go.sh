#!/bin/sh

OUT_DIR=.
PACK_NAME="models"
SUFFIX_NAME=""
READ_ONLY="created_at, updated_at"
TABLE_NAME=""
WITH_OUT=""
TAGS="bson"
TINYINT_TO_BOOL="deleted,is_admin,disable"
DSN_URL="mysql://root:12345678@127.0.0.1:3306/test?charset=utf8"
JSON_PROPERTIES="omitempty"
SPEC_TYPES="users.extra_data=struct{}, users.is_deleted=bool"
IMPORT_MODELS="github.com/civet148/db2go/models"
#指定其他orm的标签和值(以空格分隔)
COMMON_TAGS="id=gorm:\"primarykey\" create_time=gorm:\"autoCreateTime;type:timestamp\" update_time=gorm:\"autoUpdateTime;type:timestamp\""

go build -ldflags "-s -w"

if [ $? -eq 0 ]; then
./db2go --debug --url "$DSN_URL" --out "$OUT_DIR" --table "$TABLE_NAME" --json-properties "$JSON_PROPERTIES" --enable-decimal  --spec-type "$SPEC_TYPES" \
--suffix "$SUFFIX_NAME" --package "$PACK_NAME" --readonly "$READ_ONLY" --without "$WITH_OUT" --dao dao --tinyint-as-bool "$TINYINT_TO_BOOL" \
--tag "$TAGS" --import-models $IMPORT_MODELS --common-tags "$COMMON_TAGS"

echo "generate go file ok, formatting..."
gofmt -w $OUT_DIR/$PACK_NAME

else
  echo "error: db2go build failed"
fi
