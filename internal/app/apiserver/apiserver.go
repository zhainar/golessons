package apiserver

import (
	"database/sql"
	"github.com/gorilla/sessions"
	"github.com/zhainar/awesomeProject/internal/app/store/sqlstore"
	"net/http"
)

func Start(config *Config) error {
	db, err := newDB(config.DatabaseURL)

	if err != nil {
		return err
	}

	defer db.Close()

	store := sqlstore.New(db)
	sessionsStore := sessions.NewCookieStore([]byte(config.SessionsKey))

	server := newServer(store, sessionsStore)

	return http.ListenAndServe(config.BindAddr, server)
}

func newDB(databaseURL string) (*sql.DB, error) {
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
