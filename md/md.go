/*
	The MIT License

	Copyright (c) 2020 John Siu

	Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

	The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

	THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package md

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"sync"

	"github.com/J-Siu/go-helper"
	"github.com/J-Siu/go-hugo-lc/site"
)

// MD - Markdown structure
type MD struct {
	Fh    *os.File
	Links [][][]byte // all links
	Fails [][][]byte // all failed links
	Path  string
	Buf   []byte // content buffer
}

// ChkExt - check external
var ChkExt = false

// ChkWeb - check again website
var ChkWeb = false

// mds - MD array
var mds = []*MD{}

// wg - wait group
var wg sync.WaitGroup

// linkReg match [*](*)
var linkReg = regexp.MustCompile(`(\[[^[]*\])\(([^(]*)\)`)

func walkdir(path string, info os.FileInfo, err error) error {
	if info != nil {
		if !info.IsDir() && strings.ToLower(filepath.Ext(path)) == ".md" {
			m := new(MD)
			m.Path = path
			mds = append(mds, m)
			wg.Add(1)
			go m.process(&wg)
		}
	}
	return nil
}

// Process - create MD array entry
func Process(dir string) {
	helper.DebugLog("MD:Init:dir:", dir)
	// Get MD file list
	helper.ErrCheck(filepath.Walk(dir, walkdir))
	wg.Wait()
}

// Report - print
func Report() {
	var totalLink = 0
	var totalFail = 0
	for _, m := range mds {
		fmt.Printf("File: %s\n", m.Path)
		fmt.Printf("Link: %d\n", len(m.Links))
		totalLink += len(m.Links)
		totalFail += len(m.Fails)
		if m.Fails != nil {
			fmt.Printf("Link: %d\n", len(m.Fails))
			for _, fail := range m.Fails {
				fmt.Println("[x]", string(fail[2][:]))
			}
			fmt.Println()
			fmt.Println("---")
		}
	}
	fmt.Printf("Total File: %d\n", len(mds))
	fmt.Printf("Total Link: %d\n", totalLink)
	fmt.Printf("Total Fail: %d\n", totalFail)
}

// CheckFile - Check against local file
func (md *MD) checkLink(wg *sync.WaitGroup, link [][]byte) {
	var localPath string
	linkURLprep := string(link[2][:])

	if strings.HasPrefix(linkURLprep, "//") {
		linkURLprep = "https:" + linkURLprep
	}
	helper.DebugLog("MD:checkLink:linkURLprep:", linkURLprep)

	linkURL, e := url.Parse(linkURLprep)
	helper.ErrCheck(e)
	helper.DebugLog("MD:checkLink:linkURL.Host:", linkURL.Host)
	helper.DebugLog("MD:checkLink:linkURL.Path:", linkURL.Path)

	if linkURL.Host == "" {
		helper.DebugLog("MD:checkLink:(local)")
		// check if public+path exist
		localPath = path.Join(site.Site.Public, linkURL.Path)
		_, e = os.Stat(localPath)
		if e == nil {
			helper.DebugLog("MD:checkLink:localPath:(found)", localPath)
		} else {
			// path does not exist
			helper.DebugLog("MD:checkLink:localPath:(not found)", localPath)
			md.Fails = append(md.Fails, link)
		}
	} else {
		if ChkExt {
			resp, e := http.Get(linkURLprep)
			if e == nil {
				defer resp.Body.Close()
				helper.DebugLog("MD:checkLink:resp.StatusCode:", resp.StatusCode)
				if resp.StatusCode >= 400 {
					md.Fails = append(md.Fails, link)
				}
			} else {
				helper.DebugLog("MD:checkLink:ChkExt:e:", e)
				md.Fails = append(md.Fails, link)
			}
		} else {
			helper.DebugLog("MD:checkLink:(not local)")
		}
	}

	wg.Done()
}

// Check - check internal links
func (md *MD) chk() {
	// Get links
	md.Links = linkReg.FindAllSubmatch([]byte(md.Buf), -1)
	helper.DebugLog("MD:Check:md.Links#:", len(md.Links))
	// free the buf
	md.Buf = nil

	var wg sync.WaitGroup
	for _, link := range md.Links {
		wg.Add(1)
		go md.checkLink(&wg, link)
	}
	wg.Wait()
}

// Close markdown file
func (md *MD) close() error {
	helper.DebugLog("MD:Close")
	return md.Fh.Close()
}

// Open markdown file
func (md *MD) open() error {
	var e error
	helper.DebugLog("MD:Open")
	md.Fh, e = os.Open(md.Path)
	return e
}

// Read markdown file
func (md *MD) read() error {
	helper.DebugLog("MD:Read")

	var e error
	var n int
	var n64 int64

	// Reset to file start
	n64, e = md.Fh.Seek(0, 0)
	helper.ErrCheck(e)
	helper.DebugLog("MD:Read:Seek:", n64)

	stat, e := os.Stat(md.Path)
	helper.ErrCheck(e)
	md.Buf = make([]byte, stat.Size())
	n, e = md.Fh.Read(md.Buf)
	helper.DebugLog("MD:Read:byte:", n)

	return e
}

// Process markdown file
func (md *MD) process(wg *sync.WaitGroup) {
	helper.ErrCheck(md.open())
	helper.ErrCheck(md.read())
	helper.ErrCheck(md.close())
	md.chk()
	wg.Done()
}
