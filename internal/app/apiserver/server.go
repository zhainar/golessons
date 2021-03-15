package apiserver

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/sirupsen/logrus"
	"github.com/zhainar/awesomeProject/internal/app/model"
	"github.com/zhainar/awesomeProject/internal/app/store"
	"net/http"
	"time"
)

const (
	sessionName        = "sessions_data"
	ctxUserKey  ctxKey = iota
	ctxRequestID
)

var (
	errIncorrectEmailOrPassword = errors.New("incorrect email or password")
	errNotAuthenticated         = errors.New("not authenticated")
)

type server struct {
	router        *mux.Router
	logger        *logrus.Logger
	store         store.Store
	sessionsStore sessions.Store
}

type ctxKey uint8

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func newServer(store store.Store, sessionsStore sessions.Store) *server {
	s := &server{
		router:        mux.NewRouter(),
		logger:        logrus.New(),
		store:         store,
		sessionsStore: sessionsStore,
	}

	s.configureRouter()

	return s
}

func (s *server) error(w http.ResponseWriter, r *http.Request, code int, err error) {
	s.respond(w, r, code, map[string]string{"error": err.Error()})
}

func (s *server) respond(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
	w.WriteHeader(code)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}

func (s *server) configureRouter() {
	s.router.Use(s.setRequestID)
	s.router.Use(s.logRequest)
	s.router.Use(handlers.CORS(handlers.AllowedOrigins([]string{"*"})))
	s.router.HandleFunc("/users", s.handleUsersCreate()).Methods(http.MethodPost)
	s.router.HandleFunc("/sessions", s.handleSessionsCreate()).Methods(http.MethodPost)

	private := s.router.PathPrefix("/private").Subrouter()
	private.Use(s.authenticateUser)
	private.HandleFunc("/whoami", s.handlerWhoAmI()).Methods(http.MethodGet)
}

func (s *server) handleUsersCreate() http.HandlerFunc {
	type request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}

		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		u := &model.User{
			Email:    req.Email,
			Password: req.Password,
		}

		if err := s.store.User().Create(u); err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		u.Sanitize()
		s.respond(w, r, http.StatusCreated, u)
	}
}

func (s *server) handlerWhoAmI() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		s.respond(writer, request, http.StatusOK, request.Context().Value(ctxUserKey).(*model.User))
	}
}

func (s *server) handleSessionsCreate() http.HandlerFunc {
	type request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}

		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		u, err := s.store.User().FindByEmail(req.Email)

		if err != nil || !u.ComparePassword(req.Password) {
			s.error(w, r, http.StatusUnauthorized, errIncorrectEmailOrPassword)
			return
		}

		session, err := s.sessionsStore.Get(r, sessionName)

		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		session.Values["user_id"] = u.ID
		err = s.sessionsStore.Save(r, w, session)

		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		s.respond(w, r, http.StatusOK, nil)
	}
}

func (s *server) authenticateUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		session, err := s.sessionsStore.Get(request, sessionName)

		if err != nil {
			s.error(writer, request, http.StatusInternalServerError, err)
			return
		}

		id, ok := session.Values["user_id"]

		if !ok {
			s.error(writer, request, http.StatusUnauthorized, errNotAuthenticated)
			return
		}

		u, err := s.store.User().Find(id.(int))

		if err != nil {
			s.error(writer, request, http.StatusUnauthorized, errNotAuthenticated)
			return
		}

		newContext := context.WithValue(request.Context(), ctxUserKey, u)

		next.ServeHTTP(writer, request.WithContext(newContext))
	})
}

func (s *server) setRequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		id := uuid.New().String()
		writer.Header().Set("X-Request-ID", id)

		newContext := context.WithValue(request.Context(), ctxRequestID, id)

		next.ServeHTTP(writer, request.WithContext(newContext))
	})
}

func (s *server) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		logger := s.logger.WithFields(logrus.Fields{
			"remote_addr": request.RemoteAddr,
			"request_id":  request.Context().Value(ctxRequestID),
		})

		logger.Infof("started %s %s", request.Method, request.RequestURI)

		start := time.Now()
		rw := &responseWriter{writer, http.StatusOK}

		next.ServeHTTP(rw, request)

		logger.Infof(
			"completed with %d %s in %v",
			rw.code,
			http.StatusText(rw.code),
			time.Now().Sub(start),
		)
	})
}
