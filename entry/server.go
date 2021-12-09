package entry

import (
	"errors"
	"fileSync/bar"
	"fileSync/core"
	"net"
)

// 接收客户端的Header响应
// func (fe *FileEntry) RecvHeaderResponse(conn net.Conn) error {
// 	buf := make([]byte, 8)
// 	readN, err := conn.Read(buf)
// 	if err != nil {
// 		return err
// 	}
// 	recvCheckSum := core.BytesToInt64(buf[:readN])
// 	// 检查校验和
// 	if recvCheckSum == -1 {
// 		// 跳过该文件
// 		fe.FileSize = 0
// 		return nil
// 	}
// 	if fe.CheckSum != recvCheckSum {
// 		errMsg := fmt.Sprintf("检查校验和失败, 校验和不一致: %s; Bytes: %v; Header校验和: %s\n", buf[:readN], fe.CheckSum, buf[:readN])
// 		return errors.New(errMsg)
// 	}
// 	return nil
// }

// 发送文件本体
func (fe *FileEntry) Send(conn net.Conn) error {

	buf := make([]byte, 4096)
	totalSend := 0
	var bar bar.Bar
	bar.NewOption(int64(totalSend), fe.header.FileSize)
	defer bar.Finish()
	defer fe.file.Close()

	for totalSend < int(fe.header.FileSize) {
		readN, err := fe.file.Read(buf)
		if err != nil {
			return err
		}
		_, err = conn.Write(buf[:readN])
		if err != nil {
			return err
		}

		totalSend += readN
		bar.Play(int64(totalSend))

		if totalSend > int(fe.header.FileSize) {
			return errors.New("文件发送出错，发送总量大于文件")
		}
	}
	bar.Finish()
	return nil
}

// 接收客户端文件传输完成的校验
// 接收客户端接收到的文件大小，与自身的文件大小做判断
func (fe *FileEntry) Finish(conn net.Conn) error {

	buf := make([]byte, 8)
	readN, err := conn.Read(buf)
	if err != nil {
		return err
	}
	finishCheck := core.BytesToInt64(buf[:readN])

	if fe.header.FileSize != finishCheck {
		return errors.New("文件传输后校验出错")
	}
	return nil
}
