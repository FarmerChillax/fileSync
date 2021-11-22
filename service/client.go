package service

import (
	"fileSync/core"
	"fileSync/entry"
	"fmt"
	"io"
	"log"
	"net"
)

func Client(host, port string) {
	address := net.JoinHostPort(host, port)
	fmt.Printf("准备建立于 %v的链接\n", address)
	conn, _ := net.Dial("tcp", address)
	fmt.Println("链接建立成功")
	defer conn.Close()

	buf := make([]byte, 4096)
	for {
		// 获取 Header Entry
		r, err := conn.Read(buf)
		if err != nil {
			if err == io.EOF {
				fmt.Println("文件传输完成...")
			}
			return
		}
		fe := entry.FileEntry{}
		err = core.StructDecode(buf[:r], &fe)
		if err != nil {
			log.Println(err)
		}

		fmt.Printf("开始接收文件: %v\n", fe)
		// 发送响应
		err = fe.ResponseHeader(conn)
		if err != nil {
			log.Printf("响应服务端失败, err: %v\n", err)
			return
		}
		// 接收文件本体
		totalRecv, err := fe.RecvFile(conn)
		if err != nil {
			log.Printf("接收文件出错, err: %v\n", err)
			return
		}
		// 接收完成, 校验
		err = fe.Close(conn, int64(totalRecv))
		if err != nil {
			log.Printf("文件接收后校验出错, err: %v\n", err)
			return
		}
		fmt.Printf("成功接收文件: %v\n", fe)
	}

}
