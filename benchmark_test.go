package main

import (
	"os"
	"testing"
)

func BenchmarkParsingIPv4(b *testing.B) { //nolint:deadcode
	b.ReportAllocs()

	var argset argumentSet

	b.ResetTimer()

	for index := 0; index < b.N; index++ {
		argset = parseArgs([]string{os.Args[0], "1.1.1.1"})
	}
	b.StopTimer()

	result := argset

	if result.inputIP == nil {
		b.Fail()
	}
}

func BenchmarkParsingIPv6(b *testing.B) { //nolint:deadcode
	b.ReportAllocs()

	var argset argumentSet

	b.ResetTimer()

	for index := 0; index < b.N; index++ {
		argset = parseArgs([]string{os.Args[0], "2606:4700:4700::1111"})
	}
	b.StopTimer()

	result := argset

	if result.inputIP == nil {
		b.Fail()
	}
}

func BenchmarkInitialDownload(b *testing.B) { //nolint:deadcode
	b.ReportAllocs()

	for index := 0; index < b.N; index++ {
		getBlocklistFilesFromSource()
	}
}

func BenchmarkAfterInitialDownload(b *testing.B) { //nolint:deadcode
	b.ReportAllocs()
	parseArgs([]string{os.Args[0], "1.1.1.1"})
	b.ResetTimer()

	for index := 0; index < b.N; index++ {
		getBlocklistFilesFromSource()
	}
}
