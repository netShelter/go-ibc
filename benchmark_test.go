package main

import (
	"net"
	"os"
	"testing"
)

func BenchmarkParsing(b *testing.B) { //nolint:deadcode
	b.ReportAllocs()
	var inputIP net.IP
	b.ResetTimer()
	for index := 0; index < b.N; index++ {
		inputIP = parseArgs([]string{os.Args[0], "1.1.1.1"})
	}
	b.StopTimer()
	result := inputIP
	if result == nil {
		b.Fail()
	}
}

func BenchmarkInitialDownload(b *testing.B) { //nolint:deadcode
	b.ReportAllocs()
	parseArgs([]string{os.Args[0], "--ipv4", "1.1.1.1"})
	b.ResetTimer()
	for index := 0; index < b.N; index++ {
		initDownload()
	}
}

func BenchmarkAfterInitialDownload(b *testing.B) { //nolint:deadcode
	b.ReportAllocs()
	parseArgs([]string{os.Args[0], "--ipv4", "1.1.1.1"})
	b.ResetTimer()
	for index := 0; index < b.N; index++ {
		initDownload()
	}
}

/*
func BenchmarkCheckIPParsing(b *testing.B) { //nolint:deadcode
	b.ReportAllocs()
	var input ipset
	parseArgs(&input, []string{os.Args[0], "--ipv4", "1.1.1.1"})
	dir := initDownload()
	for index := 0; index < b.N; index++ {
		checkIPs(dir, &input)
	}
}
*/
