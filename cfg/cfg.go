package cfg

import (
	"os"

	"github.com/wonderivan/logger"
	"gopkg.in/yaml.v2"
)

type ConfigInfo struct {
	Server struct {
		Port string `yaml:"port"`
	}
	Db struct {
		Host          string `yaml:"host"`
		Port          string `yaml:"port"`
		Password      string `yaml:"password"`
		Db            string `yaml:"db"`
		User          string `yaml:"user"`
		MaxConnection int    `yaml:"max_connection"`
	}
}

func readFile(fileName string) []byte {
	s, err := os.ReadFile(fileName)
	if err != nil {
		logger.Error("Read file %s content error[%s]", fileName, err)
		return nil
	}

	return s
}

func LoadCfg() *ConfigInfo {
	var cfg ConfigInfo
	content := readFile("config.yaml")
	if content == nil {
		return nil
	}
	err := yaml.Unmarshal(content, &cfg)
	if err != nil {
		logger.Error("parse config.yaml error: %v", err)
		return nil
	}
	return &cfg
}
