package cfg

import (
	"os"
	"strings"

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

	envs := map[string]string{}
	envStrings := os.Environ()
	for _, env := range envStrings {
		kv := strings.Split(env, "=")
		if len(kv) == 2 {
			envs[kv[0]] = kv[1]
		}
	}

	address, ok := envs["MYSQL_ADDRESS"]
	if ok {
		kv := strings.Split(address, ":")
		cfg.Db.Host = kv[0]
		cfg.Db.Port = kv[1]
	}

	passwd, ok := envs["MYSQL_PASSWORD"]
	if ok {
		cfg.Db.Password = passwd
	}

	username, ok := envs["MYSQL_USERNAME"]
	if ok {
		cfg.Db.User = username
	}

	dbname, ok := envs["MYSQL_DB"]
	if ok {
		cfg.Db.Db = dbname
	}
	return &cfg
}
