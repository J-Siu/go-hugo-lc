/*
Copyright © 2025 John, Sing Dao, Siu <john.sd.siu@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/

package cmd

import (
	"errors"
	"net/url"
	"os"

	"github.com/J-Siu/go-helper/v2/errs"
	"github.com/J-Siu/go-helper/v2/ezlog"
	"github.com/J-Siu/go-helper/v2/ver"
	"github.com/J-Siu/go-hugo-lc/global"
	"github.com/J-Siu/go-hugo-lc/md"
	"github.com/J-Siu/go-hugo-lc/site"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "go-hugo-lc",
	Short:   `Hugo site link check`,
	Version: ver.String(),
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		ezlog.SetLogLevel(ezlog.ERR)
		if global.Flag.Debug {
			ezlog.SetLogLevel(ezlog.DEBUG)
		}
		ezlog.Debug().N("Version").Mn(ver.String()).Nn("Flag").M(&global.Flag).Out()
		// Pre-check
		if global.Flag.BaseURL == "" || site.Content == "" || site.Public == "" {
			errs.Queue("", errors.New("Missing parameters."))
		}
		if errs.IsEmpty() {
			var e error
			site.BaseURL, e = url.Parse(global.Flag.BaseURL)
			errs.Queue("", e)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		if errs.IsEmpty() {
			ezlog.Debug().
				N("BaseURL.host").Mn(site.BaseURL.Host).
				N("BaseURL.path").Mn(site.BaseURL.Path).
				N("Content").Mn(site.Content).
				N("Public").M(site.Public).
				Out()
			if global.Flag.Debug {
				ezlog.SetLogLevel(ezlog.DEBUG)
			}
			md.Process(global.Flag.Debug)
			md.Report()
		}
	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		if !errs.IsEmpty() {
			ezlog.Err().Nn("Error").M(errs.Errs).Out()
			os.Exit(1)
		}
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&global.Flag.Debug, "debug", "d", false, "Enable debug")
	rootCmd.PersistentFlags().StringVarP(&global.Flag.BaseURL, "baseURL", "b", "", "Base URL")
	rootCmd.PersistentFlags().StringVarP(&site.Content, "content", "c", "", "Content directory")
	rootCmd.PersistentFlags().StringVarP(&site.Public, "public", "p", "", "Public directory")
}
