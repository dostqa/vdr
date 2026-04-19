package config

import (
	"fmt"
	"log"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type HTTPServer struct {
	Address     string        `yaml:"address"`
	Timeout     time.Duration `yaml:"timeout"`
	IdleTimeout time.Duration `yaml:"idle_timeout"`
}

type Database struct {
	Name           string `yaml:"name"`
	Host           string `yaml:"host"`
	Port           string `yaml:"port"`
	Username       string `yaml:"username"`
	Password       string `yaml:"password"`
	MigrationsPath string `yaml:"migration_path"`
}

type FileStorage struct {
	Address    string `yaml:"address"`
	Username   string `yaml:"username"`
	Password   string `yaml:"password"`
	BucketName string `yaml:"bucketname"`
}

type Kafka struct {
	Address string `yaml:"address"`
}

type Config struct {
	Env         string      `yaml:"env"`
	Database    Database    `yaml:"database"`
	FileStorage FileStorage `yaml:"filestorage"`
	Kafka       Kafka       `yaml:"kafka"`
	HTTPServer  HTTPServer  `yaml:"http_server"`
}

func (cfg Config) DatabaseURL() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s", cfg.Database.Username, cfg.Database.Password, cfg.Database.Host, cfg.Database.Port, cfg.Database.Name)
}

func MustLoad(filename string) Config {
	var cfg Config

	if err := cleanenv.ReadConfig(fmt.Sprintf("./configs/%s", filename), &cfg); err != nil {
		log.Fatal("cannot read config: " + err.Error())
	}

	return cfg
}
