package config

import (
	"log"
	"net/url"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

// Logger
type Logger struct {
	Logfile   string `yaml:"logfile"`
	LevelInfo string `yaml:"levelInfo"`
}

// NatsConfig - NATS configuration settings
type NatsConfig struct {
	URL      string
	Host     string   `yaml:"host"`
	Subjects []string `yaml:"subjects"`
}

func (nats *NatsConfig) SetURI() {
	natsURL := &url.URL{
		Scheme: "nats",
		Host:   nats.Host,
		// User:   url.UserPassword(nats.User, nats.Password),
	}

	nats.URL = natsURL.String()
}

// Config is a configuration struct that store environmental variables
type Config struct {
	Addr             string      `yaml:"Addr"`
	ListenerProtocol string      `yaml:"listenerProtocol"`
	Logger           *Logger     `yaml:"logger"`
	Nats             *NatsConfig `yaml:"nats"`
	PostgreSQL       *PostgresConfig
	Docker           *DockerConfig
}

func (cfg *Config) GetAlt() {
	cfg.Addr = cfg.Docker.Hosts.ListenerHost

	cfg.Nats.Host = cfg.Docker.Hosts.NatsHost
	cfg.PostgreSQL.Host = cfg.Docker.Hosts.PostgresqlHost

	cfg.Nats.SetURI()
	cfg.PostgreSQL.SetURI()
}

func Setup() *Config {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("CONFIG_PATH is not set")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist: %s", configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("can't read config: %s", err)
	}

	cfg.Nats.SetURI()

	cfg.PostgreSQL = NewDB()

	if os.Getenv("DOCKER_PATH") != "" {
		cfg.Docker = GetDocker()
		cfg.GetAlt()
	}

	return &cfg
}
