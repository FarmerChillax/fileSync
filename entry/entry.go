package entry

import (
	"bytes"
	"errors"
	"fileSync/core"
	"math/rand"
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
	// control type
	// 0: push; 1: pull;
	Type         uint8
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
		file:     f,
	}, nil
}

func (fe *FileEntry) GetHeader() Header {
	return *fe.header
}

func (fe *FileEntry) GetFileName() string {
	return string(fe.filename)
}
