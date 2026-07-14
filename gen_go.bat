@echo off

rem 输出文件根目录
set OUT_DIR=.
rem 数据模型文件包名
set PACK_NAME="models"
rem 只读字段(不更新)
set READ_ONLY="created_at, updated_at"
rem 指定或排除表名(不指定则整个数据库全部导出, 排除表名在表名前面加-)
set TABLE_NAME="-user_roles"
rem 忽略字段名(逗号分隔)
set WITH_OUT=""
rem 添加标签
set TAGS="gorm"
rem 数据库连接源DSN
set DSN_URL="mysql://root:123456@127.0.0.1:3306/test?charset=utf8"
rem 指定具体表对应字段类型(不指定表则全局生效)
set SPEC_TYPES="users.extra_data=struct%%, users.is_deleted=bool"
rem 数据库DDL文件
set DDL_FILE="deploy/test.sql"


rem 判断本地系统是否已安装db2go工具，没有则进行安装
where db2go.exe

IF "%errorlevel%" == "0" (
    echo db2go already installed.
) ELSE (
    echo db2go not found in system %%PATH%%, installing...
    go install github.com/civet148/db2go@latest
    If "%errorlevel%" == "0" (
        echo db2go install succeeded
    ) ELSE (
        rem 安装失败，Linux/Mac请安装gcc工具链，Windows系统可以通过链接直接下载二进制(https://github.com/civet148/release/tree/master/db2go)
        echo error: db2go install failed, Linux/Mac please install gcc tool-chain, Windows download from https://github.com/civet148/release/tree/master/db2go
    )
)


db2go --url "%DSN_URL%" --out "%OUT_DIR%" --table "%TABLE_NAME%" --enable-decimal  --spec-type "%SPEC_TYPES%" \
 --package "%PACK_NAME%" --readonly "%READ_ONLY%" --without "%WITH_OUT%" --tag "%TAGS%" --ddl "%DDL_FILE%"

echo "generate go file ok, formatting..."
gofmt -w %OUT_DIR%/%PACK_NAME%
pause
