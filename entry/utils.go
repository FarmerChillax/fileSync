package entry

import (
	"errors"
	"fileSync/core"
	"os"
	"path/filepath"
)

// 对服务端发送的文件名进行预处理
// 传入绝对路径
func preRecvFile(p string) error {
	fPath := filepath.Dir(p)
	// 不存在则创建路径
	if !IsExist(fPath) {
		err := os.MkdirAll(fPath, 0600)
		if err != nil {
			return errors.New("创建路径失败")
		}
	}
	return nil
}

func GetExistFileSize(filename string) int64 {
	if IsExist(filename) {
		info, err := os.Stat(filename)
		if core.HandleError("检查是否跳过该文件", err) {
			return 0
		}
		return info.Size()
	}
	return 0
}
func IsExist(p string) bool {
	_, err := os.Stat(p)
	return err == nil || os.IsExist(err)
}
