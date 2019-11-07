package main

import (
	"bufio"
	"net"
	"os"
	"path/filepath"
	"strings"
)

func parserIpset(scnr *bufio.Scanner, argset argumentSet, file *os.File) (match listEntry) {
	match.match = false
	match.release = false

	for scnr.Scan() {
		switch {
		case strings.HasPrefix(scnr.Text(), "#"):
			if strings.HasPrefix(scnr.Text(), "# Maintainer  ") {
				match.maintainer = strings.TrimSpace(strings.ReplaceAll(strings.TrimPrefix(scnr.Text(), "# Maintainer"), ":", ""))
			}

			if strings.HasPrefix(scnr.Text(), "# List source URL") {
				match.url = strings.TrimSpace(strings.Replace(strings.TrimPrefix(scnr.Text(), "# List source URL"), ":", "", 1))
			}

			if strings.HasPrefix(scnr.Text(), "# Category") {
				match.category = strings.TrimSpace(strings.ReplaceAll(strings.TrimPrefix(scnr.Text(), "# Category"), ":", ""))
			}

			// Parse Listname
			match.list = strings.TrimSuffix(strings.TrimSuffix(filepath.Base(file.Name()), ".netset"), ".ipset")
		case scnr.Bytes()[0] == 0x04:
			break
		default:
			if parsedIP := net.ParseIP(scnr.Text()); parsedIP != nil {
				if parsedIP.Equal(argset.inputIP) {
					match.match = true
					match.ip = parsedIP.String()

					return match
				}
			}
		}
	}

	return match
}

func parserNetset(scnr *bufio.Scanner, argset argumentSet, file *os.File) (match listEntry) {
	match.match = false
	match.release = false

	for scnr.Scan() {
		switch {
		case strings.HasPrefix(scnr.Text(), "#"):
			if strings.HasPrefix(scnr.Text(), "# Maintainer  ") {
				match.maintainer = strings.TrimSpace(strings.ReplaceAll(strings.TrimPrefix(scnr.Text(), "# Maintainer"), ":", ""))
			}

			if strings.HasPrefix(scnr.Text(), "# List source URL") {
				match.url = strings.TrimSpace(strings.Replace(strings.TrimPrefix(scnr.Text(), "# List source URL"), ":", "", 1))
			}

			if strings.HasPrefix(scnr.Text(), "# Category") {
				match.category = strings.TrimSpace(strings.ReplaceAll(strings.TrimPrefix(scnr.Text(), "# Category"), ":", ""))
			}

			// Parse geolocation to country
			if match.category == "geolocation" {
				tmp0 := strings.TrimSuffix(filepath.Base(file.Name()), ".netset")
				tmp1 := strings.Split(tmp0, "_")
				match.country = strings.ToUpper(tmp1[len(tmp1)-1])
			}

			// Parse Listname
			match.list = strings.TrimSuffix(strings.TrimSuffix(filepath.Base(file.Name()), ".netset"), ".ipset")

		case scnr.Bytes()[0] == 0x04:
			break
		default:
			_, parsedNet, err := net.ParseCIDR(scnr.Text())
			if err == nil && parsedNet.Contains(argset.inputIP) {
				match.match = true
				match.ip = parsedNet.String()

				return match
			}
		}
	}

	return match
}
