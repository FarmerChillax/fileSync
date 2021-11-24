package core

import "log"

func HandleError(message string, err error) bool {
	if err != nil {
		log.Printf("[%s], err: %v\n", message, err)
		return true
	}
	return false
}
