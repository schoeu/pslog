#!/bin/sh

# linux 64bit
GOOS=linux GOARCH=amd64 go build -ldflags '-w -s' -o pslog_linux64
upx -9 pslog_linux64

# linux 32bit
GOOS=linux GOARCH=386 go build -ldflags '-w -s' -o pslog_linux32
upx -9 pslog_linux32

# windows 64bit
GOOS=windows GOARCH=amd64 go build -ldflags '-w -s' -o pslog_64.exe
upx -9 pslog_64.exe

# windows 32bit
GOOS=windows GOARCH=386 go build -ldflags '-w -s' -o pslog_32.exe
upx -9 pslog_32.exe

# Mac OS X 64bit
GOOS=darwin GOARCH=amd64 go build -ldflags '-w -s' -o pslog_mac
upx -9 pslog_mac
