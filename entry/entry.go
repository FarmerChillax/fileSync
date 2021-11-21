package entry

import (
	"errors"
	"fileSync/core"
	"fmt"
	"math/rand"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

const (
	recvPath = `D:\projectCode\GithubCodes\fileSync\recvDisk`
	sendPath = `D:\projectCode\GithubCodes\new-school-sdk`
)

type Entry interface {
	Send()
	Recv()
}

type FileEntry struct {
	FileSize int64
	Filename string
	CheckSum string
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
		CheckSum: strconv.Itoa(int(rand.Int31())),
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
	buf := make([]byte, 1024)
	readN, err := conn.Read(buf)
	if err != nil {
		return err
	}
	// 检查校验和
	if fe.CheckSum != string(buf[:readN]) {
		return errors.New("检查校验和失败, 校验和不一致: " + fmt.Sprintf(string(buf[:readN])))
	}
	return nil
}

// 客户端发送响应成功信息给服务端
func (fe *FileEntry) ResponseHeader(conn net.Conn) error {
	_, err := conn.Write([]byte(fe.CheckSum))
	if err != nil {
		return err
	}
	return nil
}

// 发送文件本体
func (fe *FileEntry) SendFile(conn net.Conn) error {
	buf := make([]byte, 4096)
	totalSend := 0
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

		if totalSend > int(fe.FileSize) {
			return errors.New("文件发送出错，发送总量大于文件")
		}
	}
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
		if totalRecv > int(fe.FileSize) {
			return totalRecv, errors.New("接收出错, 接收内容超出文件大小")
		}
	}
	return totalRecv, nil
}

// 响应接收成功
// 返回接收到的文件大小
func (fe *FileEntry) Close(conn net.Conn, recvSize int) error {
	_, err := conn.Write([]byte(strconv.Itoa(recvSize)))
	if err != nil {
		return errors.New("发送接收响应失败")
	}
	return nil
}

// 检查文件
func (fe *FileEntry) Finish(conn net.Conn) error {
	buf := make([]byte, 1024)
	readN, err := conn.Read(buf)
	recv, err := strconv.Atoi(string(buf[:readN]))
	if err != nil {
		return err
	}
	if fe.FileSize != int64(recv) {
		return errors.New("文件传输后校验出错")
	}
	return nil
}
