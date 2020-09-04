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

	"github.com/J-Siu/go-helper"
	"github.com/J-Siu/go-hugo-lc/md"
	"github.com/J-Siu/go-hugo-lc/site"
	"github.com/J-Siu/go-ver"
)

func usage() {
	fmt.Println("go-hugo-lc", ver.ToStr())
	fmt.Println("License : MIT License Copyright (c) 2020 John Siu")
	fmt.Println("Support : https://github.com/J-Siu/go-hugo-lc/issues")
	fmt.Println("Debug   : export _DEBUG=true")
	fmt.Println("Usage   : go-hugo-lc <baseURL> <content dir> <public dir>")
}

func main() {
	helper.DebugEnv()

	ver.Major = 0
	ver.Minor = 5
	ver.Patch = 4

	//md.ChkExt = true

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
	site.BaseURL, e = url.Parse(args[0])
	helper.ErrCheck(e)
	site.Content = args[1]
	site.Public = args[2]

	if helper.Debug {
		helper.DebugLog("BaseURL.host:", site.BaseURL.Host)
		helper.DebugLog("BaseURL.path:", site.BaseURL.Path)
		helper.DebugLog("Content:", site.Content)
		helper.DebugLog("Public:", site.Public)
	}

	md.Process()
	md.Report()
}
