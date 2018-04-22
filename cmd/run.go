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

// File describe command 'run' may used for join generated account to target channel from config file

package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"sync"

	"github.com/ealoshinsky/epic_happy/libs"
	"github.com/urfave/cli"
)

func execRunCmd(context *cli.Context) error {

	config := libs.LoadConfig(context.GlobalString("config-file"))

	sessionFiles, reason := ioutil.ReadDir(config.DataDir)
	if reason != nil {
		fmt.Println("[-] It is not possible to read the directory with session files")
		os.Exit(1)
	}
	fmt.Println("[+] Start join", len(sessionFiles), "accounts")
	lock := &sync.WaitGroup{}
	for sessionID := range sessionFiles {
		sessionPath := config.DataDir + "/" + sessionFiles[sessionID].Name()
		lock.Add(1)
		go func(lock *sync.WaitGroup, channels *libs.Config, sessionPath string) {
			session := libs.NewSession(config.TelegramID, config.TelegramAPI, sessionPath)
			if reason := session.ConnectToServer(); reason != nil {
				fmt.Println("[-] Could not connect to DC telegram servers. Try later")
				lock.Done()
			} // connect to telegram dc
			for channelID := range channels.Channels {
				channel := channels.Channels[channelID]
				if reason := session.FindAndJoinedToChannel(channel); reason != nil {
					fmt.Println("[-] Could not join to channel by session", sessionPath)
				}
			}
			lock.Done()
		}(lock, &config, sessionPath)
	}
	lock.Wait()
	fmt.Println("[+] All account success joined to", config.Channels)
	return nil
}

var runCommand = cli.Command{
	Usage: "Join generated account to target channel from config file",
	Name:  "run",

	Action: execRunCmd,
}
