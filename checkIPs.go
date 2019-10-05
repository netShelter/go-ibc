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

func checkIPs(basedir, dir string, ips *ipset) {
	//Make input and output channel and waitgroup
	in := make(chan string, 2000)
	out := make(chan listEntry, 2000)
	worker := sync.WaitGroup{}

	worker.Add(1)

	//Start insertworker
	inputWorker(basedir, in, &worker)

	//time.Sleep(2 * time.Second)
	//Start fileworker
	for i := 0; i < runtime.NumCPU(); i++ {
		go compareWorker(in, out, ips)
	}

	//time.Sleep(2 * time.Second)

	//Start output worker
	go outputWorker(dir, basedir, in, out, &worker)

	worker.Wait()
	close(in)
	close(out)
}

func inputWorker(dir string, in chan string, worker *sync.WaitGroup) {
	worker.Add(1)
	files, err := ioutil.ReadDir(dir)
	evalErr(err)
	defer worker.Done()
	for _, file := range files {
		if !file.IsDir() && file.Name() != "" {
			in <- filepath.Join(dir, file.Name())
		}
	}
}

func compareWorker(in chan string, out chan listEntry, ips *ipset) {
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
		err0 := file.Close()
		evalErr(err0)
		err1 := os.Remove(file.Name())
		evalErr(err1)
	}
}

func outputWorker(bsdir, dir string, in chan string, out chan listEntry, worker *sync.WaitGroup) {
	worker.Add(1)

	go releaseWorker(dir, in, out)

	for {
		output := <-out
		if output.match {
			fmt.Println("IP: " + output.ip + " found in Category: " + output.category + " in List: " + output.url)
		}
		if output.release {
			worker.Done()

			// Due to initial increment of waitgroup to block while executing workers
			worker.Done()
			err0 := os.Remove(dir)
			evalErr(err0)
			err1 := os.Remove(bsdir)
			evalErr(err1)
			return
		}
	}
}

func releaseWorker(dir string, in chan string, out chan listEntry) {
	for {
		if files, _ := ioutil.ReadDir(dir); len(out) == 0 && len(in) == 0 && len(files) == 0 {
			tmp := listEntry{}
			tmp.release = true
			out <- tmp
		}
		time.Sleep(2 * time.Second)
	}
}
