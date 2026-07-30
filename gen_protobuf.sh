#!/bin/sh

OUT_DIR=.
PACK_NAME="proto"
WITH_OUT="created_at, updated_at"
TABLE_NAME="users, classes"
#DSN_URL="mssql://sa:123456@127.0.0.1:1433/test?instance=SQLEXPRESS&windows=false"
#DSN_URL="postgres://postgres:123456@127.0.0.1:5432/test?sslmode=disable"
DSN_URL="mysql://root:123456@127.0.0.1:3306/test?charset=utf8"

make && ./db2go proto --url "${DSN_URL}" --out "${OUT_DIR}" --table "${TABLE_NAME}" \
                      --suffix "${SUFFIX_NAME}" --package "${PACK_NAME}" --without "${WITH_OUT}"
echo "generate protobuf file ok"
