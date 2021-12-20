package main

import (
	"fileSync/core"
	"fileSync/service"
	"log"
	"time"
)

func main() {

	host := core.Conf.Host
	port := core.Conf.Port
	startTime := time.Now()
	service.Client(host, port)
	log.Printf("传输耗时: %v\n", time.Since(startTime))
}
