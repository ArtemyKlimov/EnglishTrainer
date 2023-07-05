package store

import (
	"fmt"
	"os"
)

type Config struct {
	DatabaseURL string
	DbUser      string
	DbPwd       string
	DbHost      string
	DbPort      string
	DbName      string
}

func NewConfig() *Config {
	config := &Config{
		DbHost: os.Getenv("PG_HOSTNAME"),
		DbPort: os.Getenv("PG_PORT"),
		DbUser: os.Getenv("PG_USER"),
		DbPwd:  os.Getenv("PG_PASSWORD"),
		DbName: os.Getenv("PG_DBNAME"),
	}
	config.formatDatabaseURL()
	return config
}

func (c *Config) formatDatabaseURL() {
	c.DatabaseURL = fmt.Sprintf("host=%s dbname=%s sslmode=disable user=%s password=%s", c.DbHost, c.DbName, c.DbUser, c.DbPwd)
}
