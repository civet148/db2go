@echo off
go build -ldflags "-s -w"

set OUT_DIR=.
set PACK_NAME=dataobject
set SUFFIX_NAME="do"
set READ_ONLY="created_at, updated_at"
set DB_NAME="test"
set TABLE_NAME="users, classes"
set WITH_OUT=""
set DSN_URL="postgres://postgres:123456@127.0.0.1:5432/test?sslmode=disable"

If "%errorlevel%" == "0" (
db2go --url %DSN_URL% --out %OUT_DIR% --db %DB_NAME% --table %TABLE_NAME% ^
 --suffix %SUFFIX_NAME% --package %PACK_NAME% --readonly %READ_ONLY% --without %WITH_OUT%
echo generate go file ok, formatting...
gofmt -w %OUT_DIR%/%PACK_NAME%
)

pause