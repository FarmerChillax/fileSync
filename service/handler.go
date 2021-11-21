package service

import (
	"context"
	"fileSync/entry"
	"fmt"
	"io/fs"
	"log"
	"net"
	"path/filepath"
)

var root = `D:\projectCode\GithubCodes\new-school-sdk`
var maxRescoures = 2

// 处理TCP链接
// 遍历文件夹，发送文件
func HandleConn(conn net.Conn) {
	defer conn.Close()
	fmt.Println("链接--ok")
	tasksChan := make(chan *entry.FileEntry, maxRescoures)
	// finishChan := make(chan *entry.FileEntry)
	// finishTasks := []entry.FileEntry{}
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
	fmt.Println("传输完成，TCP通道关闭...")

}

func worker(conn net.Conn, tasksChan chan *entry.FileEntry,
	cancel context.CancelFunc) {
	// 退出时关闭goroutine
	defer cancel()

	for task := range tasksChan {
		fmt.Printf("worker get a task: %v\n", task)
		// 发送header entry
		sendN, err := task.SendHeader(conn)
		fmt.Println("发送Header成功，Header大小:", sendN)
		if err != nil {
			log.Println("发送Header失败, err:", err)
			return
		}
		// 接收响应信息
		err = task.RecvHeaderResponse(conn)
		if err != nil {
			log.Println("接收响应信息失败, err:", err)
			return
		}
		fmt.Println("校验和校验成功, 可以发送文件....")
		// 发送file本体
		err = task.SendFile(conn)
		if err != nil {
			log.Println("发送文件出错, err:", err)
			return
		}
		// 校验客户端接收的文件
		err = task.Finish(conn)
		if err != nil {
			log.Println("发送文件后校验出错, err:", err)
			return
		}
		fmt.Println("发送一个文件完成...")
	}
	// 关闭输出通道
	fmt.Println("传输完成, 通道关闭...")

}
