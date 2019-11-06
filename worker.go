package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"
)

func startWorker(dir string, argset argumentSet) {
	in := make(chan string, 2000)
	out := make(chan listEntry, 2000)
	workerWG := sync.WaitGroup{}

	workerWG.Add(1)

	inputWorker(dir, in, &workerWG)

	time.Sleep(2 * time.Second)

	for i := 0; i < runtime.NumCPU(); i++ {
		go compareWorker(in, out, argset)
	}

	time.Sleep(2 * time.Second)

	go outputWorker(in, out, &workerWG)

	workerWG.Wait()
	close(in)
	close(out)
}

func inputWorker(dir string, in chan string, workerWG *sync.WaitGroup) {
	workerWG.Add(1)

	files, err := ioutil.ReadDir(dir)
	evalErr(err, dir)

	defer workerWG.Done()

	for _, file := range files {
		if !file.IsDir() && file.Name() != "" {
			in <- filepath.Join(dir, file.Name())
		}
	}
}

func compareWorker(in chan string, out chan listEntry, argset argumentSet) {
	for {
		path := <-in
		file, err := os.Open(path)
		evalErr(err, path)

		scnr := bufio.NewScanner(file)

		if strings.HasSuffix(path, "ipset") {
			out <- scannerIpset(scnr, argset, file)
		}

		if strings.HasSuffix(path, "netset") {
			out <- scannerNetset(scnr, argset, file)
		}

		err0 := file.Close()
		evalErr(err0, file.Name())
	}
}

func outputWorker(in chan string, out chan listEntry, workerWG *sync.WaitGroup) {
	workerWG.Add(1)

	go releaseWorker(in, out)

	for {
		output := <-out
		if output.match {
			switch output.category {
			case "geolocation":
				fmt.Println("IP:" + output.ip + " | List: " + output.list + " | Country: " +
					output.country + " | URL: " + output.url)
			default:
				fmt.Println("IP:" + output.ip + " | List: " + output.list + " | Category: " +
					output.category + " | URL: " + output.url)
			}
		}

		if output.release {
			workerWG.Done()

			// Due to initial increment of waitgroup to block while executing workers
			workerWG.Done()

			return
		}
	}
}

func releaseWorker(in chan string, out chan listEntry) {
	for {
		if len(out) == 0 && len(in) == 0 {
			tmp := listEntry{}
			tmp.release = true
			out <- tmp
		}

		time.Sleep(2 * time.Second)
	}
}
