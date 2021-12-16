package entry

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"net"
)

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

// 发送文件名
func (fe *FileEntry) SendFileName(conn net.Conn) error {
	_, err := conn.Write(fe.filename)
	return err
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
