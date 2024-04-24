#!/bin/sh
go env -w GOSUMDB=off
go env -w GOPROXY=https://goproxy.cn,direct
go mod tidy