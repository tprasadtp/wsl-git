#!/bin/bash

# notice how we avoid spaces in $now to avoid quotation hell in go build
$now = Get-Date -UFormat "%Y-%m-%d_%T"
$sha1 = (git rev-parse HEAD).Trim()
$env:GOARCH = "amd64"
go build -ldflags -o .\wsl-git.exe