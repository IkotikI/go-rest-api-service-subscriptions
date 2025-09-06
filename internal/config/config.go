package config

import (
	"os"
	"strings"
	"time"

	"log"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env        string     `yaml:"env" env-default:"local"`
	DB         DB         `yaml:"db" env-required:"true"`
	HTTPServer HTTPServer `yaml:"http_server" env-required:"true"`
}

type DB struct {
	Driver       string `yaml:"driver" env-required:"true"`
	MigrationDir string `yaml:"migration_dir" env-default:"schema"`
	Username     string `yaml:"username" env-required:"true"`
	Password     string `yaml:"password" env:"POSTGRES_PASSWORD" env-required:"true"`
	Host         string `yaml:"host" env-required:"true"`
	Port         string `yaml:"port" env-required:"true"`
	DBname       string `yaml:"dbname" env-required:"true"`
	SSLMode      string `yaml:"sslmode" env-default:"disable"`
}

type HTTPServer struct {
	Debug       bool          `yaml:"debug" env-default:"false"`
	Auth        bool          `yaml:"auth" env-default:"false"`
	Addr        string        `yaml:"addr" env-default:"localhost:8080"`
	Path        string        `yaml:"path" env-default:"/"`
	Timeout     time.Duration `yaml:"timeout" env-default:"4s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"30s"`
	Users       []string      `yaml:"users"`
}

func (s HTTPServer) GetUsers() map[string]string {
	users := make(map[string]string, len(s.Users))
	for _, user := range s.Users {
		pair := strings.Split(user, ":")
		users[pair[0]] = pair[1]
	}
	return users
}

// type ServiceUser struct {
// 	User     string `yaml:"login" env-required:"true"`
// 	Password string `yaml:"password" env-required:"true"`
// }

func MustLoad() *Config {

	// log.Printf("POSTGRES_PASSWORD: %s", os.Getenv("POSTGRES_PASSWORD"))

	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("CONFIG_PATH is not set")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file %s does not exist", configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}

	if err := cleanenv.ReadEnv(&cfg); err != nil {
		log.Fatalf("cannot read env: %s", err)
	}

	return &cfg
}
