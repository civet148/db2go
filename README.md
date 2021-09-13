# db2go is a command to export database table structure to go or proto file 

## Usage

--url       connection url of database [required]

--db        databases to export, eg. "test" [required]

--table     tables to export, eg. "users, devices" [optional]

--out       output directory, default . [optional]

--package   export to a directory which used to be a package name in golang [optional]

--prefix    prefix of go file name [optional]

--suffix    suffix of go file name [optional]

--tag       customer golang struct member tag [optional] 

--without   ignore columns [optional]

--readonly  specify columns only for select [optional]

--proto     generate .proto file [optional]

--gogo-options specify gogo proto generate options [optional]

--one-file  all table schema integrate into one file named by database [optional]

--enable-decimal decimal as sqlca.Decimal type when exporting [optional]

--orm generate ORM code inner data object [optional]

--omitempty omit empty for json tag [optional]

--json-properties custom properties for json [optional]

## 1. 数据库表导出到go文件

* Windows batch 脚本

```batch
@echo off
set OUT_DIR=.
set PACK_NAME=dataobject
set SUFFIX_NAME=do
set READ_ONLY="created_at, updated_at"
set DB_NAME="test"
set TABLE_NAME="users, classes"
set WITH_OUT=""
set DSN_URL="mysql://root:123456@127.0.0.1:3306/test?charset=utf8"

db2go.exe --url %DSN_URL%  ^
--out %OUT_DIR% --db %DB_NAME% --table %TABLE_NAME% --suffix %SUFFIX_NAME% --package %PACK_NAME% --readonly %READ_ONLY% --without %WITH_OUT%

echo generate go file ok, formatting...
gofmt -w %OUT_DIR%/%PACK_NAME%
pause
```


## 2. 数据库表导出到proto文件

```batch
@echo off
set OUT_DIR=.
set PACK_NAME=proto
set WITH_OUT="created_at, updated_at"
set GOGO_OPTIONS="(gogoproto.marshaler_all)=true,(gogoproto.sizer_all)=true,(gogoproto.unmarshaler_all)=true,(gogoproto.gostring_all)=true"
set DB_NAME="test"
set TABLE_NAME="users, classes"
set SUFFIX_NAME="do"

db2go.exe --url "mysql://root:123456@127.0.0.1:3306/test?charset=utf8"  --proto --gogo-options %GOGO_OPTIONS% ^
--out %OUT_DIR% --db %DB_NAME% --table %TABLE_NAME% --suffix %SUFFIX_NAME% --package %PACK_NAME% --one-file --without %WITH_OUT%

echo generate protobuf file ok
pause
```
