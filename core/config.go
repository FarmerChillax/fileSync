package core

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"sync"
)

type Config struct {
	Host     string `json:host`
	Port     string `json:port`
	SyncRoot string `json:syncRoot`
}

var (
	Conf *Config
	once sync.Once
)

func init() {
	once.Do(func() {
		conf, err := InitConfig("./config.json")
		if err != nil {
			HandleError("初始化配置文件", err)
			return
		}
		// conf := &Config{}
		Conf = conf
	})
}

func (c *Config) Init(host, port, syncRoot string) {
	c.Host = host
	c.Port = port
	c.SyncRoot = syncRoot
}

// 读取配置文件
// 传入文件路径，e.g ./config.json
func InitConfig(filename string) (*Config, error) {
	jsonBytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var conf Config
	err = json.Unmarshal(jsonBytes, &conf)
	if err != nil {
		return nil, errors.New("解析配置文件失败")
	}
	return &conf, nil
}

// type ServiceConfig struct {
// 	Config
// }

// type ClientConfig struct {
// 	Config
// }
