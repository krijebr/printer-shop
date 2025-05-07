package config

import (
	"encoding/json"
	"io"
	"os"
)

type (
	Postgres struct {
		Host     string `json:"host"`
		Port     int    `json:"port"`
		UserName string `json:"user_name"`
		Password string `json:"password"`
		DBName   string `json:"db_name"`
	}
	HttpServer struct {
		Port int `json:"port"`
	}
	Config struct {
		Postgres   Postgres   `json:"postgres"`
		HttpServer HttpServer `json:"http_server"`
	}
)

func InitConfigFromJson(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}
	config := Config{}
	err = json.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}
