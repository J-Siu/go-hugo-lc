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
	"github.com/J-Siu/go-helper/v2/file"
	"github.com/J-Siu/go-hugo-lc/global"
	"github.com/J-Siu/go-hugo-lc/md"
	"github.com/J-Siu/go-hugo-lc/site"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "go-hugo-lc",
	Short:   `Hugo site link check`,
	Version: global.Version,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if global.Flag.Debug {
			ezlog.SetLogLevel(ezlog.DEBUG)
		}
		ezlog.Debug().N("Version").Mn(global.Version).Nn("Flag").Mn(&global.Flag).
			N("Content").Mn(site.Content).
			N("Public").M(site.Public).
			Out()

		// Pre-check

		if global.Flag.BaseURL == "" || site.Content == "" || site.Public == "" {
			cmd.Usage()
			os.Exit(0) // just exit
		}
		if errs.IsEmpty() && !file.IsDir(site.Content) {
			errs.Queue("", errors.New("Not directory: "+site.Content))
		}
		if errs.IsEmpty() && !file.IsDir(site.Public) {
			errs.Queue("", errors.New("Not directory: "+site.Public))
		}
		if errs.IsEmpty() {
			var e error
			site.BaseURL, e = url.Parse(global.Flag.BaseURL)
			errs.Queue("", e)
		}
		if errs.IsEmpty() {
			ezlog.Debug().
				N("BaseURL.host").Mn(site.BaseURL.Host).
				N("BaseURL.path").Mn(site.BaseURL.Path).
				N("Content").Mn(site.Content).
				N("Public").M(site.Public).
				Out()
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		if errs.IsEmpty() {
			if global.Flag.Debug {
				ezlog.SetLogLevel(ezlog.DEBUG)
			}
			md.Process(global.Flag.Debug)
			md.Report()
		}
	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		if errs.NotEmpty() {
			ezlog.Err().L().M(errs.Errs).Out()
			cmd.Usage()
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
	rootCmd.PersistentFlags().StringVarP(&global.Flag.BaseURL, "baseURL", "b", "", "(required) Base URL")
	rootCmd.PersistentFlags().StringVarP(&site.Content, "content", "c", "", "(required) Content directory")
	rootCmd.PersistentFlags().StringVarP(&site.Public, "public", "p", "", "(required) Public directory")
}
