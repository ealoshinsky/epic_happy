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
	"github.com/urfave/cli"
	"errors"
	"github.com/ealoshinsky/epic_happy/libs"
	"fmt"
)

// exec check subargs and prepare some data before start work
func exec(context *cli.Context) error {

	var (
		backend 		= context.String("backend")
		countNumbers 	= context.Int("count")
	)

	// load configuration file
	config := libs.LoadConfig(context.GlobalString("config-file"))
	fmt.Println(config)

	if backend == "" {
		return errors.New("[-] Missing required parameter: no backend specified")
	}
	if countNumbers == 0 {
		return errors.New("[-] Missing required parameter: count of phone number not specified")
	}

	return nil
}

var generateCommand = cli.Command{
	Usage: "Generate virtual sim(phone number) from passed backend",
	Name: "generate",
	Aliases: []string{"g"},
	Action:  exec, // <- exec command function,
	// describer sub args
	Flags: []cli.Flag{
		cli.StringFlag{
			// backend used for get phone number
			Usage: "backend used for get phone number",
			Name: "backend, b",
			Value: "",
		},
		cli.IntFlag{
			// how many number need
			Usage: "how many number needs",
			Name: "count,c",
			Value: 0,
		},
	},
}