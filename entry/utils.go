package entry

import (
	"errors"
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

func IsExist(p string) bool {
	_, err := os.Stat(p)
	return err == nil || os.IsExist(err)
}
