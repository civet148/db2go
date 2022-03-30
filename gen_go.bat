@echo off
go build -ldflags "-s -w"

set OUT_DIR=.
set PACK_NAME="models"
set SUFFIX_NAME="do"
set READ_ONLY="created_at, updated_at"
set DB_NAME="test"
set TABLE_NAME="users, classes"
set WITH_OUT=""
set TINYINT_TO_BOOL="deleted,is_admin"
rem set DSN_URL="mssql://sa:123456@127.0.0.1:1433/test?instance=SQLEXPRESS&windows=false"
rem set DSN_URL="postgres://postgres:123456@127.0.0.1:5432/test?sslmode=disable"
set DSN_URL="mysql://root:123456@127.0.0.1:3306/test?charset=utf8"
SET JSON_PROPERTIES="omitempty"

If "%errorlevel%" == "0" (
db2go --url %DSN_URL% --out %OUT_DIR% --db %DB_NAME% --table %TABLE_NAME% --json-properties %JSON_PROPERTIES% --enable-decimal ^
--suffix %SUFFIX_NAME% --package %PACK_NAME% --readonly %READ_ONLY% --without %WITH_OUT% --dao dao --tinyint-as-bool %TINYINT_TO_BOOL% ^
--import-models "github.com/civet148/db2go/models"

echo generate go file ok, formatting...
gofmt -w %OUT_DIR%/%PACK_NAME%
)
pause