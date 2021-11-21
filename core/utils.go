package core

import (
	"bytes"
	"crypto/md5"
	"encoding/gob"
	"io"
	"log"
	"os"
	"unsafe"
)

// func NewRea
func StructToByte(data interface{}) []byte {
	return *(*[]byte)(unsafe.Pointer(&data))
}

func StructEncode(data interface{}) ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	enc := gob.NewEncoder(buf)
	err := enc.Encode(data)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func StructDecode(b []byte, target interface{}) error {
	buf := bytes.NewBuffer(b)
	dec := gob.NewDecoder(buf)
	return dec.Decode(target)
}

func GetFileMD5(path string) []byte {
	f, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	h := md5.New()
	if _, err := io.Copy(h, f); err != nil {
		log.Fatal(err)
	}

	return h.Sum(nil)
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}

	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
