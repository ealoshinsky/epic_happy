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

// File describe command 'generate' may used for getting virtual sim (phone number) from other service

package cmd

import (
	"errors"
	"strings"

	"github.com/ealoshinsky/epic_happy/libs"
	"github.com/urfave/cli"

	"fmt"

	api "github.com/ealoshinsky/epic_happy/libs/backend"
)

// exec check subargs and prepare some data before start work
func exec(context *cli.Context) error {

	var (
		backend      = context.String("backend")
		countNumbers = context.Int("count")
	)

	// load configuration file
	config := libs.LoadConfig(context.GlobalString("config-file"))
	if config.DataDir == "" {
		fmt.Println("[!] Use default configuration file path and data storage")
		config.DataDir = libs.GetHomeDirectory() + "/epdata"
	} else if !strings.HasPrefix(config.DataDir, "/") {
		fmt.Println("[!] Use current directory for data save")
	}

	if backend == "" {
		return errors.New("[-] Missing required parameter: no backend specified")
	}
	if countNumbers == 0 {
		return errors.New("[-] Missing required parameter: count of phone number not specified")
	}

	// let's choice used backend
	if strings.ToLower(backend) == "simsms" {
		api.ExecuteSimSmsOrg(countNumbers, &config)
	} else if strings.ToLower(backend) == "smsko" {
		api.ExecuteSmskoRu(countNumbers, &config)
	} else {
		return errors.New("[-] Could not resolve backend for generate account")
	}

	return nil
}

var generateCommand = cli.Command{
	Usage:   "Generate virtual sim(phone number) from passed backend",
	Name:    "generate",
	Aliases: []string{"g"},
	Action:  exec, // <- exec command function,
	// describer sub args
	Flags: []cli.Flag{
		cli.StringFlag{
			// backend used for get phone number
			Usage: "backend used for get phone number",
			Name:  "backend, b",
			Value: "",
		},
		cli.IntFlag{
			// how many number need
			Usage: "how many number needs",
			Name:  "count,c",
			Value: 0,
		},
	},
}
