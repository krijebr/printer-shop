package config

import (
	"encoding/json"
	"io"
	"log/slog"
	"os"
	"time"
)

type (
	Duration time.Duration

	//RoleConf map[path]map[method][]roles
	RoleConf map[string]map[string][]string
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
	Redis struct {
		Host     string `json:"host"`
		Port     int    `json:"port"`
		Password string `json:"password"`
		DB       int    `json:"db"`
	}
	Security struct {
		TokenTTL        Duration `json:"token_ttl"`
		RefreshTokenTTL Duration `json:"refresh_token_ttl"`
		HashSalt        string   `json:"hash_salt"`
	}

	Logging struct {
		Level slog.Level `json:"level"`
	}

	Config struct {
		Postgres   Postgres   `json:"postgres"`
		HttpServer HttpServer `json:"http_server"`
		Redis      Redis      `json:"redis"`
		Security   Security   `json:"security"`
		Logging    Logging    `json:"logging"`
	}
)

func (d *Duration) UnmarshalJSON(b []byte) (err error) {
	var str string

	if err := json.Unmarshal(b, &str); err != nil {
		return err
	}
	dur, err := time.ParseDuration(str)
	*d = Duration(dur)
	return err
}

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

func InitRoleConfigFromJson(path string) (*RoleConf, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}
	roleConfig := make(RoleConf)
	err = json.Unmarshal(data, &roleConfig)
	if err != nil {
		return nil, err
	}
	return &roleConfig, nil
}
