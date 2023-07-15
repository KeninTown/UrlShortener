package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Env         string `yaml:"env" env-default:"local"`
	StoragePath string `yaml:"storage_path"`
	HttpServer  `yaml:"http_server"`
}

type HttpServer struct {
	Address string `yaml:"address"`
}

func MustLoad() *Config {
	var config Config

	if _, err := os.Stat("../configs/local.yaml"); os.IsNotExist(err) {
		log.Fatal("File is not exist")
	}

	configData, err := os.ReadFile("../configs/local.yaml")

	if err != nil {
		log.Fatalf("Can not read file: %s", err)
	}

	if err = yaml.Unmarshal(configData, &config); err != nil {
		log.Fatalf("Can not decode config file: %s", err)
	}

	return &config
}

// env: "local"
// storage_path: "path_to_storage"

// http_server:
//   address: "localhost:80"
