package core

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

type Config struct {
	Ip       string `json:"ip"`
	Port     string `json:"port"`
	Path     string `json:"path"`
	SyncTime string `json:"synctime"`
	Addr     string
}

var (
	config *Config
	once   sync.Once
)

func init() {
	config = GetConfig()
}

func GetConfig() *Config {
	if config == nil {
		once.Do(func() {
			config = &Config{}
			pwd, err := filepath.Abs(filepath.Dir(os.Args[0]))
			if err != nil {
				fmt.Println("Pwd err:", err)
				os.Exit(1)
			}
			f, err := os.Open(filepath.ToSlash(pwd + "/config.json"))
			if err != nil {
				os.Exit(1)
			}
			if err = json.NewDecoder(f).Decode(config); err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}
		})
	}
	return config
}
