package config

import (
	"log"
	"net/url"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type PostgresConfig struct {
	URL      string
	Host     string `yaml:"host"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	NameDB   string `yaml:"namedb"`
}

func NewDB() *PostgresConfig {
	pathDB := os.Getenv("DB_PATH")
	if pathDB == "" {
		log.Fatal("DB_PATH is not set")
	}

	if _, err := os.Stat(pathDB); os.IsNotExist(err) {
		log.Fatalf("config file does not exist: %s", pathDB)
	}

	var db PostgresConfig

	if err := cleanenv.ReadConfig(pathDB, &db); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}

	db.SetURI()

	return &db
}

func (db *PostgresConfig) SetURI() {
	pg := url.URL{
		Scheme: "postgres",
		Host:   db.Host,
		User:   url.UserPassword(db.User, db.Password),
		Path:   db.NameDB,
	}
	query := pg.Query()
	query.Add("sslmode", "disable")

	pg.RawQuery = query.Encode()

	db.URL = pg.String()
}
