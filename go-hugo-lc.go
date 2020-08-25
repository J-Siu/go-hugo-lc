/*
	The MIT License

	Copyright (c) 2020 John Siu

	Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

	The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

	THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package main

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/J-Siu/go-hugo-lc/helper"
	"github.com/J-Siu/go-hugo-lc/md"
	"github.com/J-Siu/go-hugo-lc/site"
)

func usage() {
	fmt.Println("go-hugo-lc version 0.5")
	fmt.Println("License : MIT License Copyright (c) 2020 John Siu")
	fmt.Println("Support : https://github.com/J-Siu/go-hugo-lc/issues")
	fmt.Println("Debug   : export _DEBUG=true")
	fmt.Println("Usage   : go-hugo-lc <baseURL> <content dir> <public dir>")
}

func main() {
	if os.Getenv("_DEBUG") == "true" {
		helper.Debug = true
	}

	// ARGs
	args := os.Args[1:]
	argc := len(args)
	switch {
	case argc == 0:
		usage()
		os.Exit(0)
	case argc == 1:
		helper.ErrCheck(errors.New("Content dir missing"))
	case argc > 3:
		usage()
		os.Exit(1)
	}
	var e error
	site.Site.BaseURL, e = url.Parse(args[0])
	helper.ErrCheck(e)
	site.Site.Content = args[1]
	site.Site.Public = args[2]
	mds := []*md.MD{}

	helper.DebugLog("BaseURL.host:", site.Site.BaseURL.Host)
	helper.DebugLog("BaseURL.path:", site.Site.BaseURL.Path)
	helper.DebugLog("Content:", site.Site.Content)
	helper.DebugLog("Public:", site.Site.Public)

	// Get MD file list
	helper.ErrCheck(filepath.Walk(site.Site.Content, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() && strings.ToLower(filepath.Ext(path)) == ".md" {
			m := new(md.MD)
			m.Path = path
			mds = append(mds, m)
		}
		return nil
	}))

	for _, m := range mds {
		helper.DebugLog("m.Path:", m.Path)
		helper.ErrCheck(m.Open())
		helper.ErrCheck(m.Read())
		helper.ErrCheck(m.Close())
		m.Check()
		if helper.Debug {
			if m.Fails == nil {
				helper.DebugLog("m.Links:nil")
			} else {
				helper.DebugLog("m.Links:")
				for _, link := range m.Links {
					helper.DebugLog(string(link[2][:]))
				}
			}
		}
		if m.Fails == nil {
			helper.DebugLog("m.Fails:nil")
		} else {
			helper.DebugLog("m.Fails:---")
			fmt.Printf("File: %s\n\n", m.Path)
			for _, fail := range m.Fails {
				fmt.Println("[x]", string(fail[2][:]))
			}
			fmt.Println()
		}
		helper.DebugLog("---")
	}
}
