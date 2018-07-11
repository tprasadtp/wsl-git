#!/bin/bash

# notice how we avoid spaces in $now to avoid quotation hell in go build
$now = Get-Date -UFormat "%Y-%m-%d_%T"
$sha1 = (git rev-parse HEAD).Trim()

GOARCH=amd64 go build -ldflags "-X main.sha1ver=$sha1 -X main.buildTime=$now" -o .\wsl-git.exe