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
)

func checkIPs(basedir, dir string, ips *ipset) {
	//Make input and output channel and waitgroup
	in := make(chan string, 100)
	out := make(chan ListEntry, 100)
	worker := sync.WaitGroup{}
	compworker := sync.WaitGroup{}

	worker.Add(1)

	//Start insertworker
	go inputWorker(basedir, in, &worker)

	//Start fileworker
	for i := 0; i < runtime.NumCPU(); i++ {
		go compareWorker(in, out, ips, &worker, &compworker)
	}

	//Start output worker
	go outputWorker(out, &worker, &compworker)

	worker.Wait()
	close(in)
	close(out)
	//os.RemoveAll(dir)
}

func inputWorker(dir string, in chan string, worker *sync.WaitGroup) {
	worker.Add(1)
	files, err := ioutil.ReadDir(dir)
	evalErr(err)
	defer worker.Done()
	for _, file := range files {
		if !file.IsDir() {
			in <- filepath.Join(dir, file.Name())
		}
	}
}

func compareWorker(in chan string, out chan ListEntry, ips *ipset, worker *sync.WaitGroup, compworker *sync.WaitGroup) {
	worker.Add(1)
	compworker.Add(1)
	defer worker.Done()
	defer compworker.Done()
	for {
		path := <-in
		file, err := os.Open(path)
		evalErr(err)
		scnr := bufio.NewScanner(file)
		if strings.HasSuffix(path, "ipset") {
			out <- scannerIpset(scnr, ips)
		}
		if strings.HasSuffix(path, "netset") {
			out <- scannerNetset(scnr, ips)
		}
		file.Close()
	}
}

func outputWorker(out chan ListEntry, worker *sync.WaitGroup, compworker *sync.WaitGroup) {
	worker.Add(1)
	defer worker.Done()

	// Due to initial increment of waitgroup to block while executing workers
	defer worker.Done()

	go releaseWorker(out, compworker)

	for {
		output := <-out
		if output.match {
			fmt.Println("IP: " + output.ip + " found in Category: " + output.category + " in List: " + output.url)
		}
		if output.release {
			return
			//break
		}
	}
}

func releaseWorker(out chan ListEntry, compworker *sync.WaitGroup) {
	compworker.Wait()
	tmp := ListEntry{}
	tmp.release = true
	out <- tmp
}
