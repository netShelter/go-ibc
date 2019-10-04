package main

import (
	"bufio"
	"log"
	"net"
	"strings"
)

func evalErr(err error) {
	if err != nil {
		log.Fatalln("Error:", err)
	}
}

/*
func evalErrSoft(err error) {
	if err != nil {
		log.Println("Error:", err)
	}
}
*/

type ListEntry struct {
	maintainer string
	url        string
	category   string
	ip         string
	match      bool
	release    bool
}

func scannerIpset(scnr *bufio.Scanner, ips *ipset) (match ListEntry) {
	match.match = false
	match.release = false
	for scnr.Scan() {
		switch {
		case strings.HasPrefix(scnr.Text(), "#"):
			if strings.HasPrefix(scnr.Text(), "# Maintainer  ") {
				match.maintainer = strings.TrimSpace(strings.ReplaceAll(strings.TrimPrefix(scnr.Text(), "# Maintainer"), ":", ""))
			}
			if strings.HasPrefix(scnr.Text(), "# List source URL") {
				match.url = strings.TrimSpace(strings.ReplaceAll(strings.TrimPrefix(scnr.Text(), "# List source URL"), ":", ""))
			}
			if strings.HasPrefix(scnr.Text(), "# Category") {
				match.category = strings.TrimSpace(strings.ReplaceAll(strings.TrimPrefix(scnr.Text(), "# Category"), ":", ""))
			}
		case scnr.Bytes()[0] == 0x04:
			break
		default:
			if parsedIP := net.ParseIP(scnr.Text()); parsedIP != nil {
				if parsedIP.Equal(ips.ipv4) || parsedIP.Equal(ips.ipv6) {
					match.match = true
					match.ip = parsedIP.String()
					return
				}
			}
		}
	}
	return
}

func scannerNetset(scnr *bufio.Scanner, ips *ipset) (match ListEntry) {
	match.match = false
	match.release = false
	for scnr.Scan() {
		switch {
		case strings.HasPrefix(scnr.Text(), "#"):
			if strings.HasPrefix(scnr.Text(), "# Maintainer  ") {
				match.maintainer = strings.TrimSpace(strings.ReplaceAll(strings.TrimPrefix(scnr.Text(), "# Maintainer"), ":", ""))
			}
			if strings.HasPrefix(scnr.Text(), "# List source URL") {
				match.url = strings.TrimSpace(strings.ReplaceAll(strings.TrimPrefix(scnr.Text(), "# List source URL"), ":", ""))
			}
			if strings.HasPrefix(scnr.Text(), "# Category") {
				match.category = strings.TrimSpace(strings.ReplaceAll(strings.TrimPrefix(scnr.Text(), "# Category"), ":", ""))
			}
		case scnr.Bytes()[0] == 0x04:
			break
		default:
			_, parsedNet, err := net.ParseCIDR(scnr.Text())
			if err == nil && (parsedNet.Contains(ips.ipv4) || parsedNet.Contains(ips.ipv6)) {
				match.match = true
				match.ip = parsedNet.String()
				return
			}
		}
	}
	return
}
