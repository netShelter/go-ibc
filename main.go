package main

import (
	"os"
)

func main() {
	input := parseArgs(os.Args)
	dir := initDownload()
	checkIPs(dir, input)
}
