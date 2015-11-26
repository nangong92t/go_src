# Copyright 2014, Successfulmatch Inc. All rights reserved.
# Author TonyXu<tonycbcd@gmail.com>,
# Build on dev-0.0.1
# MIT Licensed

MAKEFLAGS = -s

sql.go: sql.y
	go tool yacc -o sql.go sql.y
	gofmt -w sql.go

clean:
	rm -f y.output sql.go
