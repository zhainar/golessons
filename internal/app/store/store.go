package store

import (
	"bytes"
	"database/sql"
	"errors"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
)
import _ "github.com/lib/pq"

type Store struct {
	config         *Config
	db             *sql.DB
	userRepository *UserRepository
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
		return encodeError(err)
	}

	s.db = db

	return nil
}

func (s *Store) Close() {
	s.db.Close()
}

func encodeError(err error) error {
	// --- Encoding: convert s from UTF-8 to ShiftJIS
	// declare a bytes.Buffer b and an encoder which will write into this buffer
	var b bytes.Buffer

	wInUTF8 := transform.NewWriter(&b, charmap.Windows1251.NewDecoder())

	// encode our string
	wInUTF8.Write([]byte(err.Error()))
	wInUTF8.Close()

	return errors.New(b.String())
}

func (s *Store) User() *UserRepository {
	if s.userRepository == nil {
		s.userRepository = &UserRepository{
			store: s,
		}
	}

	return s.userRepository
}