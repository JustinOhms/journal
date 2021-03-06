// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

// Command bake reads a set of files and writes a Go source file to "static.go"
// that declares a map of string constants containing contents of the input files.
package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"unicode/utf8"
)

var files = []string{
	"case_0/2014-01-01-0000-EST",
	"case_0/2014-01-02-0000-EST",
	"case_0/2014-01-03-0000-EST",
	"case_0/2014-01-04-0000-EST",
	"case_0/2014-01-05-0000-EST",
	"case_0/2014-01-06-0000-EST",
	"case_0/2014-01-07-0000-EST",
	"case_0.json",
	"case_0_fix_reflog/0",
	"case_0_fix_reflog/1",
	"case_0_fix_reflog/10",
	"case_0_fix_reflog/11",
	"case_0_fix_reflog/12",
	"case_0_fix_reflog/13",
	"case_0_fix_reflog/14",
	"case_0_fix_reflog/2",
	"case_0_fix_reflog/3",
	"case_0_fix_reflog/4",
	"case_0_fix_reflog/5",
	"case_0_fix_reflog/6",
	"case_0_fix_reflog/7",
	"case_0_fix_reflog/8",
	"case_0_fix_reflog/9",
}

func main() {
	if err := bake(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func bake() error {
	f, err := os.Create("static.go")
	if err != nil {
		return err
	}
	defer f.Close()
	w := bufio.NewWriter(f)
	fmt.Fprintf(w, "%v\n\npackage case_0_static\n\n", warning)
	fmt.Fprintf(w, "var Files = map[string]string{\n")
	for _, fn := range files {
		b, err := ioutil.ReadFile(fn)
		if err != nil {
			return err
		}
		fmt.Fprintf(w, "\t%q: ", fn)
		if utf8.Valid(b) {
			fmt.Fprintf(w, "`%s`", sanitize(b))
		} else {
			fmt.Fprintf(w, "%q", b)
		}
		fmt.Fprintln(w, ",\n")
	}
	fmt.Fprintln(w, "}")
	if err := w.Flush(); err != nil {
		return err
	}
	return f.Close()
}

// sanitize prepares a valid UTF-8 string as a raw string constant.
func sanitize(b []byte) []byte {
	// Replace ` with `+"`"+`
	b = bytes.Replace(b, []byte("`"), []byte("`+\"`\"+`"), -1)

	// Replace BOM with `+"\xEF\xBB\xBF"+`
	// (A BOM is valid UTF-8 but not permitted in Go source files.
	// I wouldn't bother handling this, but for some insane reason
	// jquery.js has a BOM somewhere in the middle.)
	return bytes.Replace(b, []byte("\xEF\xBB\xBF"), []byte("`+\"\\xEF\\xBB\\xBF\"+`"), -1)
}

const warning = "// DO NOT EDIT ** This file was generated with the bake tool ** DO NOT EDIT //"
