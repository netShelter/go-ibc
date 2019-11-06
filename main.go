package main

import (
	"os"
)

func main() {
	argset := parseArgs(os.Args)
	dir := initDownload()
	checkIPs(dir, argset)
}
