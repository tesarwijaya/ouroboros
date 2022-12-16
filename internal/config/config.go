package config

import (
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Port string `envconfig:"PORT" default:"8000"`

	SqlDBHost     string `envconfig:"APP_SQL_DB_HOST" default:"ouroboros-sql-db"`
	SqlDBPort     int64  `envconfig:"APP_SQL_DB_PORT" default:"5432"`
	SqlDBUsername string `envconfig:"APP_SQL_DB_USERNAME" default:"root"`
	SqlDBPassword string `envconfig:"APP_SQL_DB_PASSWORD" default:"pass"`
	SqlDBName     string `envconfig:"APP_SQL_DB_NAME" default:"ouroboros_db"`

	EventStoreDBHost string `envconfig:"APP_EVENT_STORE_DB_HOST" default:"eventstoredb"`
	EventStoreDBPort int64  `envconfig:"APP_EVENT_STORE_DB_PORT" default:"1113"`
}

func NewConfig() (*Config, error) {
	var c Config

	_ = godotenv.Load(".env")

	err := envconfig.Process("APP", &c)
	if err != nil {
		return nil, err
	}

	return &c, nil
}
