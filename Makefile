#SHELL=/usr/bin/env bash

build:
	go build -ldflags "-s -w" -o db2go

install:build
	sudo cp db2go /usr/local/bin

clean:
	rm -f db2go