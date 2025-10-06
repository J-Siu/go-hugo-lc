/*
	The MIT License

	Copyright (c) 2025 John Siu

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

	"github.com/J-Siu/go-helper/v2/basestruct"
	"github.com/J-Siu/go-helper/v2/ezlog"
	"github.com/J-Siu/go-hugo-lc/site"
)

var (
	// ChkExt - check external
	ChkExt = false

	// ChkWeb - check again website
	ChkWeb = false

	// mds - MD array
	mds = []*MD{}

	// wg - wait group
	wg sync.WaitGroup

	// linkReg match [*](*)
	linkReg = regexp.MustCompile(`(\[[^[]*\])\(([^(]*)\)`)

	logLevel ezlog.Level = ezlog.ERR
)

// MD - Markdown structure
type MD struct {
	basestruct.Base
	Fh    *os.File   `json:"fh,omitempty"`
	Links [][][]byte `json:"links,omitempty"` // all links
	Fails [][][]byte `json:"fails,omitempty"` // all failed links
	Path  string     `json:"path,omitempty"`
	Buf   []byte     `json:"buf,omitempty"` // content buffer
}

func (t *MD) New() *MD {
	t.MyType = "MD"
	t.Initialized = true
	return t
}

// CheckFile - Check against local file
func (t *MD) chkLink(wg *sync.WaitGroup, link [][]byte) {
	prefix := t.MyType + ".chkLink"
	log := ezlog.New().SetLogLevel(logLevel)
	log.Debug().N(prefix).Out()
	var (
		e         error
		linkURL   *url.URL
		localPath string
	)
	linkURLprep := string(link[2][:])

	if strings.HasPrefix(linkURLprep, "//") {
		linkURLprep = "https:" + linkURLprep
	}
	log.Debug().N(prefix).N("linkURLprep").M(linkURLprep).Out()

	linkURL, t.Err = url.Parse(linkURLprep)
	if t.Err == nil {
		log.Debug().N(prefix).N("linkURL.Host").M(linkURL.Host).Out()
		log.Debug().N(prefix).N("linkURL.Path").M(linkURL.Path).Out()

		if linkURL.Host == "" {
			log.Debug().N(prefix).M("(local)").Out()
			// check if public+path exist
			localPath = path.Join(site.Public, linkURL.Path)
			_, e = os.Stat(localPath)
			if e == nil {
				log.Debug().N(prefix).N("localPath:(found)").M(localPath).Out()
			} else {
				// path does not exist
				log.Debug().N(prefix).N("localPath:(not found)").M(localPath).Out()
				t.Fails = append(t.Fails, link)
			}
		} else {
			if ChkExt {
				resp, e := http.Get(linkURLprep)
				if e == nil {
					defer resp.Body.Close()
					log.Debug().N(prefix).N("resp.StatusCode").M(resp.StatusCode).Out()
					if resp.StatusCode >= 400 {
						t.Fails = append(t.Fails, link)
					}
				} else {
					log.Debug().N(prefix).N("ChkExt:e").M(e).Out()
					t.Fails = append(t.Fails, link)
				}
			} else {
				log.Debug().N(prefix).M("(not local)").Out()
			}
		}
	}
	wg.Done()
}

// Check - check internal links
func (t *MD) chk() {
	prefix := t.MyType + ".chk"
	log := ezlog.New().SetLogLevel(logLevel)
	log.Debug().N(prefix).Out()
	if t.Err == nil {
		// Get links
		t.Links = linkReg.FindAllSubmatch([]byte(t.Buf), -1)
		log.Debug().N(prefix).N("md.Links#").M(len(t.Links)).Out()
		// free the buf
		t.Buf = nil

		var wg sync.WaitGroup
		for _, link := range t.Links {
			wg.Add(1)
			go t.chkLink(&wg, link)
		}
		wg.Wait()
	}
}

// Close markdown file
func (t *MD) close() *MD {
	prefix := t.MyType + ".close"
	log := ezlog.New().SetLogLevel(logLevel)
	log.Debug().N(prefix).Out()
	if t.Err == nil {
		t.Err = t.Fh.Close()
	}
	return t
}

// Open markdown file
func (t *MD) open() *MD {
	prefix := t.MyType + ".open"
	log := ezlog.New().SetLogLevel(logLevel)
	log.Debug().N(prefix).Out()
	if t.Err == nil {
		t.Fh, t.Err = os.Open(t.Path)
	}
	return t
}

// Read markdown file
func (t *MD) read() *MD {
	prefix := t.MyType + ".read"
	log := ezlog.New().SetLogLevel(logLevel)
	log.Debug().N(prefix).Out()
	var (
		n int
		// n64  int64
		stat os.FileInfo
	)
	if t.Err == nil {
		// Reset to file start
		// n64, t.Err = t.Fh.Seek(0, 0)
		_, t.Err = t.Fh.Seek(0, 0)
	}
	if t.Err == nil {
		// log.Debug().N(prefix).N("Seek").M(n64).Out()
		stat, t.Err = os.Stat(t.Path)
	}
	if t.Err == nil {
		t.Buf = make([]byte, stat.Size())
		n, t.Err = t.Fh.Read(t.Buf)
		log.Debug().N(prefix).N("byte").M(n).Out()
	}
	return t
}

// Process markdown file
func (t *MD) process(wg *sync.WaitGroup) {
	t.open().read().close().chk()
	wg.Done()
}

func walkDir(path string, info os.FileInfo, err error) error {
	if info != nil {
		if !info.IsDir() && strings.ToLower(filepath.Ext(path)) == ".md" {
			m := new(MD).New()
			m.Path = path
			mds = append(mds, m)
			wg.Add(1)
			go m.process(&wg)
		}
	}
	return nil
}

// Process - create MD array entry
func Process(debug bool) {
	prefix := "md.Process"
	if debug {
		logLevel = ezlog.DEBUG
	}
	ezlog.Debug().N(prefix).M(site.Content).Out()
	// Get MD file list
	if filepath.Walk(site.Content, walkDir) != nil {
		return
	}
	wg.Wait()
}

// Report - print
func Report() {
	var totalLink = 0
	var totalFail = 0
	for _, m := range mds {
		fmt.Printf("Fail: %d/%d | Path: %s\n", len(m.Fails), len(m.Links), m.Path)
		totalLink += len(m.Links)
		totalFail += len(m.Fails)
		if m.Fails != nil {
			for _, fail := range m.Fails {
				fmt.Println("[x]", string(fail[2][:]))
			}
			fmt.Println("---")
		}
	}
	fmt.Printf("Total File: %d\n", len(mds))
	fmt.Printf("Total Link: %d\n", totalLink)
	fmt.Printf("Total Fail: %d\n", totalFail)
}
