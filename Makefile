#SHELL=/usr/bin/env bash

DATE_TIME=`date +'%Y%m%d %H:%M:%S'`
COMMIT_ID=`git rev-parse --short HEAD`

build:
	CGO_ENABLED=1 go build -ldflags "-s -w -X 'main.BuildTime=${DATE_TIME}' -X 'main.GitCommit=${COMMIT_ID}'"

install:
	CGO_ENABLED=1 go install -ldflags "-s -w -X 'main.BuildTime=${DATE_TIME}' -X 'main.GitCommit=${COMMIT_ID}'" && db2go -v

clean:
	rm -f db2go
