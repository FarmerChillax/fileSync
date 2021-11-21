package service

import "io/fs"

// Walk的回调函数
func GetLoaclSyncFile(path string, info fs.FileInfo, err error) error {
	if err != nil {
		return err
	}
	// 发送文件
	if !info.IsDir() {

	}
	return nil
}
