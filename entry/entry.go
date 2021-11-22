package entry

import (
	"errors"
	"fileSync/bar"
	"fileSync/core"
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"path/filepath"
	"time"
)

var (
	recvPath = core.Conf.SyncRoot
	sendPath = core.Conf.SyncRoot
)

type Entry interface {
	Send()
	Recv()
}

type FileEntry struct {
	FileSize int64
	Filename string
	CheckSum int64
	file     *os.File
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func New(filename string) (*FileEntry, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	fileInfo, err := f.Stat()
	if err != nil {
		return nil, err
	}

	fPath, err := filepath.Rel(sendPath, filename)
	if err != nil {
		return nil, errors.New("获取相对路径失败")
	}

	return &FileEntry{
		FileSize: fileInfo.Size(),
		Filename: fPath,
		CheckSum: rand.Int63(),
		file:     f,
	}, nil
}

// 发送文件Header到TCP流中
func (fe *FileEntry) SendHeader(conn net.Conn) (int, error) {
	// 编码后发送header
	feBytes, err := core.StructEncode(fe)
	if err != nil {
		return 0, err
	}
	// 发送Header
	writeN, err := conn.Write(feBytes)
	if err != nil {
		return 0, err
	}

	return writeN, nil
}

// 从TCP链接中读取数据
func (fe *FileEntry) RecvHeaderResponse(conn net.Conn) error {
	buf := make([]byte, 8)
	readN, err := conn.Read(buf)
	if err != nil {
		return err
	}
	recvCheckSum := core.BytesToInt64(buf[:readN])
	// 检查校验和
	if fe.CheckSum != recvCheckSum {
		errMsg := fmt.Sprintf("检查校验和失败, 校验和不一致: %s; Bytes: %v; Header校验和: %s\n", buf[:readN], fe.CheckSum, buf[:readN])
		return errors.New(errMsg)
	}
	return nil
}

// 客户端发送响应成功信息给服务端
func (fe *FileEntry) ResponseHeader(conn net.Conn) error {
	checkSumBytes := core.Int64ToBytes(fe.CheckSum)
	writeN, err := conn.Write([]byte(checkSumBytes))
	log.Printf("[写入TCP流校验和] 校验和长度: %v; 校验和内容: %d, 校验和Bytes: %v\n", writeN, fe.CheckSum, checkSumBytes)
	if err != nil {
		return err
	}
	return nil
}

// 发送文件本体
func (fe *FileEntry) SendFile(conn net.Conn) error {

	buf := make([]byte, 4096)
	totalSend := 0
	var bar bar.Bar
	bar.NewOption(int64(totalSend), fe.FileSize)
	defer bar.Finish()
	defer fe.file.Close()

	for totalSend < int(fe.FileSize) {
		readN, err := fe.file.Read(buf)
		if err != nil {
			return err
		}
		_, err = conn.Write(buf[:readN])
		if err != nil {
			return err
		}

		totalSend += readN
		bar.Play(int64(totalSend))

		if totalSend > int(fe.FileSize) {
			return errors.New("文件发送出错，发送总量大于文件")
		}
	}
	bar.Finish()
	return nil
}

// 接收文件本体
func (fe *FileEntry) RecvFile(conn net.Conn) (totalRecv int, err error) {
	fePath := filepath.Join(recvPath, fe.Filename)
	if fe.file == nil {
		err = preRecvFile(fePath)
		if err != nil {
			return totalRecv, err
		}

		fe.file, err = os.OpenFile(fePath, os.O_CREATE|os.O_WRONLY, 0600)
		if err != nil {
			return totalRecv, err
		}
	}

	// 开始接收文件
	var bar bar.Bar
	bar.NewOption(0, fe.FileSize)
	defer bar.Finish()
	defer fe.file.Close()

	buf := make([]byte, 4096)
	for totalRecv < int(fe.FileSize) {
		readN, err := conn.Read(buf)
		if err != nil {
			return totalRecv, err
		}
		_, err = fe.file.Write(buf[:readN])
		if err != nil {
			return totalRecv, err
		}
		totalRecv += readN
		bar.Play(int64(totalRecv))
		if totalRecv > int(fe.FileSize) {
			return totalRecv, errors.New("接收出错, 接收内容超出文件大小")
		}
	}
	bar.Finish()
	return totalRecv, nil
}

// 响应接收成功
// 返回接收到的文件大小
func (fe *FileEntry) Close(conn net.Conn, recvSize int64) error {

	_, err := conn.Write(core.Int64ToBytes(recvSize))
	if err != nil {
		return errors.New("发送接收响应失败")
	}
	return nil
}

// 检查文件
func (fe *FileEntry) Finish(conn net.Conn) error {
	buf := make([]byte, 8)
	readN, err := conn.Read(buf)
	if err != nil {
		return err
	}
	finishCheck := core.BytesToInt64(buf[:readN])

	if fe.FileSize != finishCheck {
		return errors.New("文件传输后校验出错")
	}
	return nil
}
