package configs

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
)

type Config struct {
	Postgres   Postgres   `yaml:"postgres"`
	Redis      Redis      `yaml:"redis"`
	HttpServer HttpServer `yaml:"http_server"`
	GrpcServer GrpcServer `yaml:"grpc_server"`
}

type HttpServer struct {
	Port      int    `yaml:"port"`
	ElkDomain string `yaml:"elk_domain"`
}

type GrpcServer struct {
	Port int `yaml:"port"`
}

type Redis struct {
	Address  string `yaml:"address"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

type Postgres struct {
	Host     string `yaml:"host"`
	Port     uint16 `yaml:"port"`
	Database string `yaml:"database"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}

func InitConfig(filePath string) (*Config, error) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal YAML config: %w", err)
	}

	return &config, nil
}
