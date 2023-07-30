package remote_command

import (
	"fmt"
	"github.com/go-yaml/yaml"
	"io/ioutil"
)

type RemoteNodeConfig struct {
	RemoteNodes []*RemoteNode `json:"remoteNodes" yaml:"remoteNodes"`
}

type RemoteNode struct {
	Host     string `json:"host" yaml:"host"`
	User     string `json:"user" yaml:"user"`
	Password string `json:"password" yaml:"password"`
	Port     string `json:"port" yaml:"port"`
}

func NewRemoteNodeConfig() *RemoteNodeConfig {
	return &RemoteNodeConfig{}
}

func loadConfigFile(path string) []byte {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Errorf("load file err: %s", err)
		return nil
	}
	return b
}

// LoadConfig 读取配置文件
func loadConfig(path string) (*RemoteNodeConfig, error) {
	c := NewRemoteNodeConfig()
	if b := loadConfigFile(path); b != nil {
		err := yaml.Unmarshal(b, c)
		if err != nil {
			fmt.Errorf("unmarshal err: %s", err)
			return nil, err
		}
		return c, err
	} else {
		return nil, fmt.Errorf("load config file error")
	}
}
