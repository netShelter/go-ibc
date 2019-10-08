package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"
)

func checkIPs(dir string, ips net.IP) {
	//Make input and output channel and waitgroup
	in := make(chan string, 2000)
	out := make(chan listEntry, 2000)
	worker := sync.WaitGroup{}

	worker.Add(1)

	//Start insertworker
	inputWorker(dir, in, &worker)

	time.Sleep(2 * time.Second)
	//Start fileworker
	for i := 0; i < runtime.NumCPU(); i++ {
		go compareWorker(in, out, ips)
	}

	time.Sleep(2 * time.Second)

	//Start output worker
	go outputWorker(in, out, &worker)

	worker.Wait()
	close(in)
	close(out)
}

func inputWorker(dir string, in chan string, worker *sync.WaitGroup) {
	worker.Add(1)
	files, err := ioutil.ReadDir(dir)
	evalErr(err, dir)
	defer worker.Done()
	for _, file := range files {
		if !file.IsDir() && file.Name() != "" {
			in <- filepath.Join(dir, file.Name())
		}
	}
}

func compareWorker(in chan string, out chan listEntry, ips net.IP) {
	for {
		path := <-in
		file, err := os.Open(path)
		evalErr(err, path)
		scnr := bufio.NewScanner(file)
		if strings.HasSuffix(path, "ipset") {
			out <- scannerIpset(scnr, ips, file)
		}
		if strings.HasSuffix(path, "netset") {
			out <- scannerNetset(scnr, ips, file)
		}
		err0 := file.Close()
		evalErr(err0, file.Name())
	}
}

func outputWorker(in chan string, out chan listEntry, worker *sync.WaitGroup) {
	worker.Add(1)
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
			worker.Done()

			// Due to initial increment of waitgroup to block while executing workers
			worker.Done()
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
