package store

import (
	"database/sql"

	_ "github.com/lib/pq"
)

type Store struct {
	config *Config
	Db     *sql.DB
}

func New(config *Config) *Store {
	return &Store{
		config: config,
	}
}

func (s *Store) Open() error {
	db, err := sql.Open("postgres", s.config.DatabaseURL)
	if err != nil {
		return err
	}
	if err := db.Ping(); err != nil {
		return err
	}
	s.Db = db
	return nil
}

func (s *Store) Close() {
	s.Db.Close()
}
