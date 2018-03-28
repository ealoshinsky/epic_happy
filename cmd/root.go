/*
This is free and unencumbered software released into the public domain.

Anyone is free to copy, modify, publish, use, compile, sell, or
distribute this software, either in source code form or as a compiled
binary, for any purpose, commercial or non-commercial, and by any
means.

In jurisdictions that recognize copyright laws, the author or authors
of this software dedicate any and all copyright interest in the
software to the public domain. We make this dedication for the benefit
of the public at large and to the detriment of our heirs and
successors. We intend this dedication to be an overt act of
relinquishment in perpetuity of all present and future rights to this
software under copyright law.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
IN NO EVENT SHALL THE AUTHORS BE LIABLE FOR ANY CLAIM, DAMAGES OR
OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE,
ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
OTHER DEALINGS IN THE SOFTWARE.

For more information, please refer to <http://unlicense.org>
*/

// File describer entrance point

package cmd

import (
	"github.com/urfave/cli"
	"fmt"
	"strings"
	"github.com/ealoshinsky/epic_happy/libs"
	"os"
)

// Run represent entrance point for the entire application
func Run(args []string, release, commit, buildTime, appName string){
	app := cli.NewApp()
	app.Name = appName
	app.Usage = "Getting up rating in telegram groups and channels"
	app.Version = release

	// showing version
	cli.VersionPrinter = func(c *cli.Context) {
		fmt.Println(appName, "-", strings.ToLower(app.Usage))
		fmt.Printf("Build version: %s\nFrom commit: %s\nBuild at: %s\n", release, commit, buildTime)
	}

	// set global args
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name: "config-file",
			Value: libs.GetHomeDirectory() + "/epdata/config.yml",
			Usage: "Load configuration from file",
		},

	}

	// set commands
	app.Commands = []cli.Command{
		generateCommand,
	}

	// start app
	if reason := app.Run(args);reason != nil {
		fmt.Println(reason)
		os.Exit(1)
	}

}