package config

import (
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type DockerConfig struct {
	Hosts *AlternativeHost
}

func GetDocker() *DockerConfig {
	var docker DockerConfig

	docker.Hosts = ParseAlternate()

	return &docker
}

type AlternativeHost struct {
	NatsHost       string `yaml:"natsHost"`
	PostgresqlHost string `yaml:"postgresqlHost"`
	ListenerHost   string `yaml:"listenerHost"`
}

func ParseAlternate() *AlternativeHost {
	var hosts AlternativeHost

	dockerPath := os.Getenv("DOCKER_PATH")
	if dockerPath == "" {
		log.Fatal("CONFIG_PATH is not set")
	}

	if _, err := os.Stat(dockerPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist: %s", dockerPath)
	}

	if err := cleanenv.ReadConfig(dockerPath, &hosts); err != nil {
		log.Fatalf("can't read config: %s", err)
	}

	return &hosts
}
