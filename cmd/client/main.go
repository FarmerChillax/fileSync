package main

import (
	"fileSync/core"
	"fileSync/service"
	"fmt"
)

func main() {
	fmt.Println("client")
	host := core.Conf.Host
	port := core.Conf.Port
	service.Client(host, port)
}
