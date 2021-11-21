package core

// func GetLoaclSyncFile(conn net.Conn, Path string) {
// 	defer conn.Close()
// 	filepath.Walk(Path, func(path string, info fs.FileInfo, err error) error {
// 		if err != nil {
// 			return err
// 		}
// 		if !info.IsDir() {
// 			// path 为绝对路径
// 			relPath, _ := filepath.Rel(config.Path, path)
// 			relPath = filepath.ToSlash(relPath)
// 			// 相对路径
// 			// fmt.Println("rel", relPath)
// 			// 构建文件块，然后发送文件
// 			fileEntry := NewFileEntry(relPath, info)
// 			err = fileEntry.send(conn, path)
// 			if err != nil {
// 				fmt.Println("发送文件时出错:", err)
// 				return err
// 			}
// 		}
// 		return nil
// 	})
// }
