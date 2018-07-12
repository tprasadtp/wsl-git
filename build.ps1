#!/bin/bash

# notice how we avoid spaces in $now to avoid quotation hell in go build

$env:GOARCH = "amd64"
$env:GOOS = "windows"
go build -ldflags -o .\wsl-git.exe