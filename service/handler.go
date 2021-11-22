package service

import (
	"context"
	"fileSync/core"
	"fileSync/entry"
	"fmt"
	"io/fs"
	"log"
	"net"
	"path/filepath"
)

var root = core.Conf.SyncRoot
var maxRescoures = 2

// 处理TCP链接
// 遍历文件夹，发送文件
func HandleConn(conn net.Conn) {
	defer conn.Close()
	fmt.Println("链接--ok")
	tasksChan := make(chan *entry.FileEntry, maxRescoures)
	ctx, cancel := context.WithCancel(context.Background())
	go worker(conn, tasksChan, cancel)
	fmt.Println("Worker Running...")
	// 遍历发送文件
	filepath.Walk(root, func(path string, info fs.FileInfo, err error) error {
		if !info.IsDir() {
			// 构建entry
			fe, err := entry.New(path)
			if err != nil {
				log.Println("New file entry err:", err)
				return nil
			}
			tasksChan <- fe
		}
		return nil
	})
	close(tasksChan)
	<-ctx.Done()
	log.Println("传输完成，TCP通道关闭...")

}

func worker(conn net.Conn, tasksChan chan *entry.FileEntry,
	cancel context.CancelFunc) {
	// 退出时关闭goroutine
	defer fmt.Println("worker退出")
	defer cancel()

	for task := range tasksChan {
		log.Printf("worker get a task: %v\n", task)
		// 发送header entry
		sendN, err := task.SendHeader(conn)
		log.Println("[发送Header成功], Header大小:", sendN)
		if err != nil {
			log.Println("[发送Header失败], err:", err)
			return
		}
		// 接收响应信息
		err = task.RecvHeaderResponse(conn)
		if err != nil {
			log.Println("[接收响应信息失败], err:", err)
			return
		}
		log.Println("[校验和校验成功], 可以发送文件")
		// 发送file本体
		err = task.SendFile(conn)
		if err != nil {
			log.Println("[发送文件出错], err:", err)
			return
		}
		// 校验客户端接收的文件
		err = task.Finish(conn)
		if err != nil {
			log.Println("[发送文件后校验出错], err:", err)
			return
		}
		log.Printf("[发送成功], 文件信息: %v\n", task)
	}
	// 关闭输出通道
	log.Println("[传输完成], 通道关闭...")

}
