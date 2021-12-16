package entry

import (
	"encoding/binary"
	"net"
	"os"
	"path/filepath"
)

// 接收文件Header
func (fe *FileEntry) RecvHeader(conn net.Conn) error {
	err := binary.Read(conn, binary.BigEndian, fe.header)
	return err
}

// 接收文件名
func (fe *FileEntry) RecvFileName(conn net.Conn) error {
	buf := make([]byte, fe.header.FileNameSize)
	readN, err := conn.Read(buf)
	if err != nil {
		return err
	}
	fe.filename = buf[:readN]
	return nil
}

func (fe *FileEntry) CheckExistFile(conn net.Conn) error {
	existSize := GetExistFileSize(string(fe.filename))
	if existSize != fe.header.FileSize {
		return nil
	}
	fe.header.FileSize = 0
	fe.header.IsSkip = true

	err := binary.Write(conn, binary.BigEndian, fe.header)
	if err != nil {
		return err
	}
	return nil
}

// 从tcp stream 读文件
func (fe *FileEntry) RecvFile(conn net.Conn) (totalRecv int, err error) {
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

	defer fe.file.Close()
	nextRecv := 4096
	if fe.header.FileSize < 4096 {
		nextRecv = int(fe.header.FileSize)
	}
	buf := make([]byte, nextRecv)
	for totalRecv < int(fe.header.FileSize) {
		// 读取内容
		readN, err := conn.Read(buf)
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
