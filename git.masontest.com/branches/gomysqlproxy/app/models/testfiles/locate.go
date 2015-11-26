// Copyright 2014, Successfulmatch Inc. All rights reserved.
// Author TonyXu<tonycbcd@gmail.com>,
// Build on dev-0.0.1
// MIT Licensed

// Package testfiles locates test files within the Vitess directory tree.
package testfiles

import (
	"path"
	"path/filepath"
)

func Locate(filename string) string {
	root := "/usr/local/go/goextentions/src/git.masontest.com/branches/gomysqlproxy/app/models"
	return path.Join(root, "data", "test", filename)
}

func Glob(pattern string) []string {
	root := "/usr/local/go/goextentions/src/git.masontest.com/branches/gomysqlproxy/app/models"
	resolved := path.Join(root, "data", "test", pattern)
	out, err := filepath.Glob(resolved)
	if err != nil {
		panic(err)
	}
	return out
}
