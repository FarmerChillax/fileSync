package core

import "log"

// 处理错误信息
// message: 消息前缀, err: error体
func HandleError(message string, err error) bool {
	if err != nil {
		log.Printf("[%s], err: %v\n", message, err)
		return true
	}
	return false
}
