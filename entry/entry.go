package entry

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"fileSync/core"
	"fmt"
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

type Header struct {
	FileSize     int64
	FileNameSize int64
}

type FileEntry struct {
	header   *Header
	filename []byte
	file     *os.File
}

func init() {
	// 初始化随机种子，用于传输校验
	rand.Seed(time.Now().UnixNano())
}

func GetEmpty() *FileEntry {
	return &FileEntry{
		header: &Header{},
	}
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

	filePathBuffer := bytes.NewBufferString(fPath)
	return &FileEntry{
		header: &Header{
			FileSize:     fileInfo.Size(),
			FileNameSize: int64(filePathBuffer.Len()),
		},
		filename: filePathBuffer.Bytes(),
		// CheckSum: rand.Int63(),
		file: f,
	}, nil
}

func (fe *FileEntry) GetHeader() Header {
	return *fe.header
}

func (fe *FileEntry) GetFileName() string {
	return string(fe.filename)
}

// 发送文件Header
func (fe *FileEntry) SendHeader(conn net.Conn) error {
	buf := &bytes.Buffer{}
	err := binary.Write(buf, binary.BigEndian, fe.header)
	if err != nil {
		return err
	}
	_, err = conn.Write(buf.Bytes())
	fmt.Printf("发送出去的Header: %v; bytes: %v\n", fe.header, buf.Bytes())
	return err
}

// 接收文件Header
func (fe *FileEntry) RecvHeader(conn net.Conn) error {
	err := binary.Read(conn, binary.BigEndian, fe.header)
	return err
}

// 发送文件名
func (fe *FileEntry) SendFileName(conn net.Conn) error {
	_, err := conn.Write(fe.filename)
	return err
}

// 接收文件名
func (fe *FileEntry) RecvFileName(conn net.Conn) error {
	buf := make([]byte, fe.header.FileNameSize)
	reader := bufio.NewReader(conn)
	_, err := reader.Read(buf)
	// filename, err := reader.Peek(int(fe.header.FileNameSize))
	if err != nil {
		return err
	}
	fe.filename = buf
	return nil
}

// 往tcp stream写文件
func (fe *FileEntry) SendFile(conn net.Conn) error {
	buf := make([]byte, 4096)
	totalSend := 0

	defer fe.file.Close()
	for totalSend < int(fe.header.FileSize) {
		readN, err := fe.file.Read(buf)
		if err != nil {
			return err
		}
		writeN, err := conn.Write(buf[:readN])
		if err != nil {
			return err
		}
		totalSend += writeN
		if totalSend > int(fe.header.FileSize) {
			return errors.New("文件发送错误, 发送总量大于文件")
		}
	}

	fmt.Println("发送结束，发送大小:", totalSend)
	return nil
}

// 从tcp stream 读文件
func (fe *FileEntry) RecvFile(conn net.Conn) (totalRecv int, err error) {
	// fmt.Println()
	// time.Sleep(time.Second * 3)

	fePath := filepath.Join(recvPath, string(fe.filename))
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
	fmt.Println("创建空间成功")
	// reader := bufio.NewReader(conn)
	// writer := bufio.NewWriter(fe.file)
	// rw := bufio.NewReadWriter(reader, writer)
	defer fe.file.Close()
	nextRecv := 4096
	if fe.header.FileSize < 4096 {
		nextRecv = int(fe.header.FileSize)
	}
	buf := make([]byte, nextRecv)
	for totalRecv < int(fe.header.FileSize) {
		// 读取内容
		readN, err := conn.Read(buf)
		// fmt.Printf("next recv: %v; readN: %v; total: %v\n", nextRecv, readN, totalRecv)
		if err != nil {
			return totalRecv, err
		}
		// 写入文件
		_, err = fe.file.Write(buf[:readN])
		if err != nil {
			return totalRecv, err
		}
		totalRecv += readN
		if fe.header.FileSize-int64(totalRecv) < int64(nextRecv) {
			nextRecv = int(fe.header.FileSize) - totalRecv
			buf = make([]byte, nextRecv)
		}
	}

	return
}
