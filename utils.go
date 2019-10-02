package main

import "log"

func evalErr(err error) {
	if err != nil {
		log.Fatalln("Error:", err)
	}
}
