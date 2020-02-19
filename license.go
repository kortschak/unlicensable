// Copyright ©2019 Dan Kortschak. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// unlicensable lists licenses within a source tree, noting possible missed licenses.
package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/licensecheck"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatalf("Usage: %s <path>", os.Args[0])
	}
	files, err := licenses(os.Args[1])
	if err != nil {
		fmt.Printf("error walking the path %q: %v\n", os.Args[1], err)
	}
	fmt.Printf("%+v\n", files)
}

type license struct {
	path  string
	cover licensecheck.Coverage
}

func licenses(path string) ([]license, error) {
	var files []license
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.IsDir() {
			if name := info.Name(); name != "." && strings.HasPrefix(name, ".") {
				return filepath.SkipDir
			}
			return nil
		}
		name := strings.ToLower(info.Name())
		if !strings.Contains(name, "license") && !strings.Contains(name, "licence") {
			return nil
		}
		b, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}
		cover, ok := licensecheck.Cover(b, licensecheck.Options{})
		if ok {
			files = append(files, license{path: path, cover: cover})
		} else {
			log.Println("missed:", path)
		}
		return nil
	})
	return files, err
}
