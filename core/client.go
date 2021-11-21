package core

import (
	"fmt"
	"net"
)

func Client(host, port string) {
	addr := net.JoinHostPort(host, port)
	fmt.Println("connect addr:", addr)
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

}
