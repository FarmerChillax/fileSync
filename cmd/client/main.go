package main

import (
	"fileSync/service"
	"fmt"
)

func main() {
	fmt.Println("client")
	host := "127.0.0.1"
	port := "5000"
	service.Client(host, port)
}
