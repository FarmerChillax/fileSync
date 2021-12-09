package service

import (
	"fileSync/core"
	"fileSync/entry"
	"fmt"
	"io"
	"net"
)

func Client(host, port string) {
	address := net.JoinHostPort(host, port)
	fmt.Printf("准备建立于 %v的链接\n", address)
	conn, _ := net.Dial("tcp", address)
	fmt.Println("链接建立成功")
	defer conn.Close()

	for {
		fe := entry.GetEmpty()
		// get Header
		fmt.Println(fe.GetHeader())
		err := fe.RecvHeader(conn)
		if err == io.EOF {
			fmt.Println("tcp通道关闭")
			return
		}
		fmt.Println(fe.GetHeader())
		if core.HandleError("接受帧头出错", err) {
			return
		}
		// recv filename
		err = fe.RecvFileName(conn)
		if core.HandleError("接受文件名出错", err) {
			return
		}
		fmt.Printf("recv file: %s\n", fe.GetFileName())
		// save file
		readN, err := fe.RecvFile(conn)
		if core.HandleError("接受文件出错", err) {
			return
		}
		fmt.Printf("接受文件成功, 大小: %v\n", readN)
	}

}
