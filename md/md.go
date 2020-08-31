/*
	The MIT License

	Copyright (c) 2020 John Siu

	Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

	The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

	THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package md

import (
	"net/url"
	"os"
	"path"
	"regexp"
	"strings"

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

// LinkReg match [*](*)
var LinkReg = regexp.MustCompile(`(\[[^[]*\])\(([^(]*)\)`)

// Check - check internal links
func (m *MD) Check() {
	var localPath string
	// Get links
	m.Links = LinkReg.FindAllSubmatch([]byte(m.Buf), -1)
	// free the buf
	m.Buf = nil
	for _, link := range m.Links {
		linkURLprep := string(link[2][:])
		helper.DebugLog("MD:Check:linkURLprep:", linkURLprep)

		if strings.HasPrefix(linkURLprep, "//") {
			helper.DebugLog("MD:Check:+https")
			linkURLprep = "https:" + linkURLprep
		}

		linkURL, e := url.Parse(linkURLprep)
		helper.ErrCheck(e)
		helper.DebugLog("MD:Check:linkURL.Host:", linkURL.Host)
		helper.DebugLog("MD:Check:linkURL.Path:", linkURL.Path)

		if linkURL.Host == "" || linkURL.Host == site.Site.BaseURL.Host {
			helper.DebugLog("MD:Check:(local)")
			// check if public+path exist
			localPath = path.Join(site.Site.Public, linkURL.Path)
			_, e = os.Stat(localPath)
			if e == nil {
				helper.DebugLog("MD:Check:localPath:(found)", localPath)
			} else {
				// path does not exist
				helper.DebugLog("MD:Check:localPath:(not found)", localPath)
				m.Fails = append(m.Fails, link)
			}
		} else {
			helper.DebugLog("MD:Check:(not local)")
		}
	}
}

// Close markdown file
func (m *MD) Close() error {
	helper.DebugLog("MD:Close")
	return m.Fh.Close()
}

// Open markdown file
func (m *MD) Open() error {
	var e error
	helper.DebugLog("MD:Open")
	m.Fh, e = os.Open(m.Path)
	return e
}

// Read markdown file
func (m *MD) Read() error {
	helper.DebugLog("MD:Read")

	var e error
	var n int
	var n64 int64

	// Reset to file start
	n64, e = m.Fh.Seek(0, 0)
	helper.ErrCheck(e)
	helper.DebugLog("MD:Read:Seek:", n64)

	stat, e := os.Stat(m.Path)
	helper.ErrCheck(e)
	m.Buf = make([]byte, stat.Size())
	n, e = m.Fh.Read(m.Buf)
	helper.DebugLog("MD:Read:byte:", n)

	return e
}
