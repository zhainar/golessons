package apiserver

import (
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/zhainar/awesomeProject/internal/app/store"
	"net/http"
)

type server struct {
	router *mux.Router
	logger *logrus.Logger
	store  store.Store
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func newServer(store store.Store) *server {
	s := &server{
		router: mux.NewRouter(),
		logger: logrus.New(),
		store:  store,
	}

	s.configureRouter()

	return s
}

func (s *server) configureRouter() {
	s.router.HandleFunc("/users", s.handleUsersCreate()).Methods(http.MethodPost)
}

func (s *server) handleUsersCreate() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {

	}
}
