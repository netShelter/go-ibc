package main

import (
	"os"
)

func main() {
	argset := parseArgs(os.Args)
	dir := getBlocklistFilesFromSource()
	processGivenData(dir, argset)
}

func processGivenData(dir string, argset argumentSet) {
	startWorker(dir, argset)
}
