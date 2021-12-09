package entry

// import (
// 	"errors"
// 	"fileSync/bar"
// 	"fileSync/core"
// 	"fmt"
// 	"log"
// 	"net"
// 	"os"
// 	"path/filepath"
// )

// // 客户端发送响应成功信息给服务端
// func (fe *FileEntry) ResponseHeader(conn net.Conn) error {
// 	// 文件已传输则跳过
// 	fePath := filepath.Join(recvPath, fe.Filename)
// 	if alreadySize := IsSkip(fePath); alreadySize == fe.FileSize {
// 		fmt.Printf("文件已传输, 跳过该文件: %v\n", fe.Filename)
// 		fe.CheckSum = -1
// 		fe.FileSize = 0
// 	}

// 	checkSumBytes := core.Int64ToBytes(fe.CheckSum)
// 	writeN, err := conn.Write([]byte(checkSumBytes))
// 	log.Printf("[写入TCP流校验和] 校验和长度: %v; 校验和内容: %d\n", writeN, fe.CheckSum)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

// // 接收文件本体
// func (fe *FileEntry) Recv(conn net.Conn) (totalRecv int, err error) {
// 	fePath := filepath.Join(recvPath, fe.Filename)
// 	if fe.file == nil {
// 		err = preRecvFile(fePath)
// 		if err != nil {
// 			return totalRecv, err
// 		}

// 		fe.file, err = os.OpenFile(fePath, os.O_CREATE|os.O_WRONLY, 0600)
// 		if err != nil {
// 			return totalRecv, err
// 		}
// 	}

// 	// 开始接收文件
// 	var bar bar.Bar
// 	bar.NewOption(0, fe.FileSize)
// 	defer bar.Finish()
// 	defer fe.file.Close()
// 	// totalRecv, err = io.Copy(fe.file, conn)

// 	// result := bytes.NewBuffer(nil)
// 	// result.Bytes()
// 	// var newBuf [1024]byte

// 	buf := make([]byte, 4096)
// 	for totalRecv < int(fe.FileSize) {
// 		readN, err := conn.Read(buf)
// 		if err != nil {
// 			return totalRecv, err
// 		}
// 		_, err = fe.file.Write(buf[:readN])
// 		if err != nil {
// 			return totalRecv, err
// 		}
// 		totalRecv += readN
// 		bar.Play(int64(totalRecv))
// 		if totalRecv > int(fe.FileSize) {
// 			return totalRecv, errors.New("接收出错, 接收内容超出文件大小")
// 		}
// 	}

// 	bar.Finish()
// 	return totalRecv, nil
// }

// // 响应接收成功
// // 返回接收到的文件大小
// func (fe *FileEntry) Close(conn net.Conn, recvSize int64) error {

// 	_, err := conn.Write(core.Int64ToBytes(recvSize))
// 	if err != nil {
// 		return errors.New("发送接收响应失败")
// 	}
// 	return nil
// }
