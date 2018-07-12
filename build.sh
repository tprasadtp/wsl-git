#!/bin/bash

# notice how we avoid spaces in $now to avoid quotation hell in go build
env GOOS=windows GOARCH=amd64 go build -ldflags  -o ./wsl-git.exe