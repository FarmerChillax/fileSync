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
	defer log.Println("worker退出")
	defer cancel()

	for task := range tasksChan {
		// 发送header entry
		err := task.SendHeader(conn)
		if core.HandleError("发送帧头出错", err) {
			return
		}
		fmt.Println(task.GetHeader())
		// 发送文件名
		err = task.SendFileName(conn)
		if core.HandleError("发送文件名出错", err) {
			return
		}
		// 发送文件本体
		fmt.Println("成功发送文件名:", task.GetFileName())
		err = task.SendFile(conn)
		if core.HandleError("发送文件本体出错", err) {
			return
		}
		break
	}
}
