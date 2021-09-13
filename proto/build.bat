@echo off

set SRC_DIR=.\
set DST_DIR=.\
set GOGOPROTO_PATH=%GOPATH%\src\github.com\gogo\protobuf\protobuf

echo generating...

protoc -I=%SRC_DIR%  --proto_path=%GOPATH%\src;%GOGOPROTO_PATH%  --gogo_out=plugins=grpc,Mgoogle/protobuf/wrappers.proto=github.com\gogo\protobuf\types,:%DST_DIR%  %SRC_DIR%\*.proto

echo generate over

pause
