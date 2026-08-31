@echo off
go build -ldflags "-s -w"

set OUT_DIR=.
set PACK_NAME="proto"
set WITHOUT=""
set TABLE_NAME="users, classes"
rem set DSN_URL="mssql://sa:123456@127.0.0.1:1433/test?instance=SQLEXPRESS&windows=false"
rem set DSN_URL="postgres://postgres:123456@127.0.0.1:5432/test?sslmode=disable"
set DSN_URL="mysql://root:123456@127.0.0.1:3306/test?charset=utf8"

If "%errorlevel%" == "0" (
.\db2go.exe proto --url %DSN_URL% --out %OUT_DIR% --db %DB_NAME% --table %TABLE_NAME% --package %PACK_NAME% --without %WITHOUT%
echo generate protobuf file ok
)
pause