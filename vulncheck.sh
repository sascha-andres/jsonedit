#!/bin/sh

go install golang.org/x/vuln/cmd/govulncheck@latest

govulncheck ./...
