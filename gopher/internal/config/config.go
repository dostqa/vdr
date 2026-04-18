package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type HTTPServer struct {
	Address     string        `yaml:"address"`
	Timeout     time.Duration `yaml:"timeout"`
	IdleTimeout time.Duration `yaml:"idle_timeout"`
}

type DataBase struct {
	Name           string `yaml:"name" env:"DB_NAME" env-default:"db"`
	Host           string `yaml:"host" env:"DB_HOST" env-default:"db"`
	Port           string `yaml:"port" env:"DB_PORT" env-default:"5432"`
	Username       string `yaml:"username" env:"DB_USER" env-default:"admin"`
	Password       string `yaml:"password" env:"DB_PASSWORD" env-default:"admin"`
	MigrationsPath string `yaml:"migration_path" env:"MIGRATION_PATH" env-default:"migrations"`
}

type FileStorage struct {
	Address    string `yaml:"address"`
	Username   string `yaml:"username"`
	Password   string `yaml:"password"`
	BucketName string `yaml:"bucketname"`
}

type Config struct {
	Env         string      `yaml:"env"`
	StoragePath string      `yaml:"storage_path"`
	DataBase    DataBase    `yaml:"database"`
	FileStorage FileStorage `yaml:"filebase"`
	HTTPServer  `yaml:"http_server"`
}

func (cfg Config) DataBaseURL() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s", cfg.DataBase.Username, cfg.DataBase.Password, cfg.DataBase.Host, cfg.DataBase.Port, cfg.DataBase.Name)
}

func NewConfigFromFile(path string) (*Config, error) {
	const op = "config.NewConfigFromFile"

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	config := &Config{}
	err = yaml.Unmarshal(data, config)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return config, nil
}
