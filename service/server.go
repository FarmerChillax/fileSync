package service

import (
	"context"
	"fmt"
	"log"
	"net"
)

// 启动服务器
func Run(ctx context.Context, host, port string) (context.Context, error) {
	address := net.JoinHostPort(host, port)
	listen, err := net.Listen("tcp", address)
	if err != nil {
		return ctx, err
	}

	ctx, cancel := context.WithCancel(ctx)

	go func() {
		for {
			conn, err := listen.Accept()
			if err != nil {
				log.Printf("Accept err: %v\n", err)
				continue
			}
			log.Printf("Accept a connect from: %s\n", conn.RemoteAddr())
			go HandleConn(conn)
		}
	}()

	go func() {
		fmt.Printf("Service is running in %v\n", address)
		fmt.Println("service started. Press any key to stop.")
		var s string
		fmt.Scanln(&s)
		cancel()
	}()

	return ctx, nil
}
