package main

import "net"

type argumentSet struct {
	inputIP net.IP
}

type listEntry struct {
	list       string
	maintainer string
	url        string
	category   string
	country    string
	ip         string
	match      bool
	release    bool
}
