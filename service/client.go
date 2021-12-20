package service

import (
	"fileSync/core"
	"fileSync/entry"
	"io"
	"log"
	"net"
)

func Client(host, port string) {
	address := net.JoinHostPort(host, port)
	log.Printf("准备建立于 %v的链接\n", address)
	conn, _ := net.Dial("tcp", address)
	defer conn.Close()

	for {
		fe := entry.GetEmpty()
		// get Header
		err := fe.RecvHeader(conn)
		if err == io.EOF {
			// 正常传输完成
			// 后期添加任务完成标记
			return
		}
		if core.HandleError("接受帧头出错", err) {
			return
		}
		// recv filename
		err = fe.RecvFileName(conn)
		if core.HandleError("接受文件名出错", err) {
			return
		}
		// check exist
		err = fe.CheckExistFile(conn)
		if core.HandleError("检测文件存在出错", err) {
			return
		}
		// save file
		readN, err := fe.RecvFile(conn)
		if core.HandleError("接受文件出错", err) {
			return
		}
		if fe.GetHeader().IsSkip {
			log.Printf("跳过传输, 文件已存在:%s", fe.GetFileName())
		} else {
			log.Printf("接受文件成功, 传输大小: %v\n", readN)
		}
	}

}
