package main

import "log"

func evalErr(err error, args ...string) {
	if err != nil {
		log.Fatalln("Error:", err, args)
	}
}
