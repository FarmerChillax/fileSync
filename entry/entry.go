package entry

import (
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

type FileEntry struct {
	FileSize int64
	Filename string
	CheckSum int64
	file     *os.File
}

func init() {
	// 初始化随机种子，用于传输校验
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
