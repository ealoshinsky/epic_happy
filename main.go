package main

import "fmt"

var (
	buildTime = "unset"
	commit = "unset"
	release = "unset"
	appName = "unset"
)

func main() {
	fmt.Println(buildTime, commit, release, appName)
}