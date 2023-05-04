@echo off
go build -ldflags "-s -w"

set OUT_DIR=.
set PACK_NAME="proto"
set WITH_OUT="created_at, updated_at"
rem set GOGO_OPTIONS="(gogoproto.marshaler_all)=true,(gogoproto.sizer_all)=true,(gogoproto.unmarshaler_all)=true,(gogoproto.gostring_all)=true"
set GOGO_OPTIONS=""
set DB_NAME="test"
set TABLE_NAME="users, classes"
set SUFFIX_NAME="do"
rem set DSN_URL="mssql://sa:123456@127.0.0.1:1433/test?instance=SQLEXPRESS&windows=false"
rem set DSN_URL="postgres://postgres:123456@127.0.0.1:5432/test?sslmode=disable"
set DSN_URL="mysql://root:123456@127.0.0.1:3306/test?charset=utf8"

If "%errorlevel%" == "0" (
.\db2go.exe --url %DSN_URL% --proto --gogo-options %GOGO_OPTIONS% ^
--out %OUT_DIR% --db %DB_NAME% --table %TABLE_NAME% --suffix %SUFFIX_NAME% --package %PACK_NAME% --one-file --without %WITH_OUT%
echo generate protobuf file ok
)
pause