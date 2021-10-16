package main

import (
	"bufio"
	"net"
	"os"
	"path/filepath"
	"strings"
)

func parseFile(scanner *bufio.Scanner, argset argumentSet, file *os.File) (match listEntry) {
	match.match = false
	match.release = false

	for scanner.Scan() {
		switch {
		case strings.HasPrefix(scanner.Text(), "#"):
			parseMetadataLine(file, scanner, &match)

		case scanner.Bytes()[0] == 0x04:
			break
		default:
			parseIPAndNet(scanner.Text(), argset, &match)
		}
	}

	return match
}

func getURL(text, token string) string {
	return strings.TrimSpace(strings.Replace(strings.TrimPrefix(text, token), ":", "", 1))
}

func getContent(text, token string) string {
	return strings.TrimSpace(strings.ReplaceAll(strings.TrimPrefix(text, token), ":", ""))
}

func parseMetadataLine(file *os.File, scanner *bufio.Scanner, match *listEntry) {
	if strings.HasPrefix(scanner.Text(), "# Maintainer  ") {
		match.maintainer = getContent(scanner.Text(), "# Maintainer")
	}

	if strings.HasPrefix(scanner.Text(), "# Maintainer URL") {
		match.url = getURL(scanner.Text(), "# Maintainer URL")
	}

	if strings.HasPrefix(scanner.Text(), "# List source URL") {
		if tmpurl := getURL(scanner.Text(), "# List source URL"); tmpurl != "" {
			match.url = tmpurl
		}
	}

	if strings.HasPrefix(scanner.Text(), "# Category") {
		match.category = getContent(scanner.Text(), "# Category")
	}

	if match.category == "geolocation" {
		pureFilename := strings.TrimSuffix(filepath.Base(file.Name()), ".netset")
		splittedFilename := strings.Split(pureFilename, "_")
		match.country = strings.ToUpper(splittedFilename[len(splittedFilename)-1])
	}

	// Parse Listname
	match.list = strings.TrimSuffix(strings.TrimSuffix(filepath.Base(file.Name()), ".netset"), ".ipset")
}

func parseIPAndNet(text string, argset argumentSet, match *listEntry) {
	_, parsedNet, err := net.ParseCIDR(text)
	if err == nil && parsedNet.Contains(argset.inputIP) {
		match.match = true
		match.ip = parsedNet.String()

		return
	}

	if parsedIP := net.ParseIP(text); parsedIP != nil {
		if parsedIP.Equal(argset.inputIP) {
			match.match = true
			match.ip = parsedIP.String()

			return
		}
	}
}
