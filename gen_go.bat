@echo off
go build -ldflags "-s -w"

set OUT_DIR=.
set PACK_NAME="models"
set SUFFIX_NAME="do"
set READ_ONLY="created_at, updated_at"
set TABLE_NAME="users,classes"
set WITH_OUT=""
set TAGS="bson"
set TINYINT_TO_BOOL="deleted,is_admin,disable"
set DSN_URL="mysql://root:123456@127.0.0.1:3306/test?charset=utf8"
set JSON_PROPERTIES="omitempty"
set SPEC_TYPES="users.extra_data=struct{}"
set IMPORT_MODELS="github.com/civet148/db2go/models"
rem 指定其他orm的标签和值(以空格分隔)
set COMMON_TAGS="id=gorm:\"primarykey\" created_at=gorm:\"autoCreateTime;type:timestamp\" updated_at=gorm:\"autoUpdateTime;type:timestamp\""

If "%errorlevel%" == "0" (
.\db2go.exe --url %DSN_URL% --out %OUT_DIR% --table %TABLE_NAME% --json-properties %JSON_PROPERTIES% --enable-decimal  --spec-type %SPEC_TYPES% ^
--suffix %SUFFIX_NAME% --package %PACK_NAME% --readonly %READ_ONLY% --without %WITH_OUT% --dao dao --tinyint-as-bool %TINYINT_TO_BOOL% ^
--tag %TAGS% --import-models %IMPORT_MODELS%

echo generate go file ok, formatting...
gofmt -w %OUT_DIR%/%PACK_NAME%
)
pause
