package article

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/google/uuid"
)

type HTTPServer struct {
	Port uint16
	Host string

	uc *ArticleUseCase

	router *chi.Mux

	articleReader ArticleReader
}

func (s HTTPServer) Host1() string {
	return s.Host
}

func (s HTTPServer) Port1() uint16 {
	return s.Port
}

func (s *HTTPServer) setupRoute() {
	r := s.router

	r.Route("/articles", func(r chi.Router) {
		r.Post("/", s.NewArticleHandler)
		r.Get("/", s.ListArticlesHandler)
		r.Delete("/", s.DeleteArticleHandler)

		r.Route("/{articleID}", func(r chi.Router) {
			r.Put("/", s.EditArticleHandler)
			r.Get("/", s.SingleArticleHandler)
			r.Delete("/", s.DeleteArticleHandler)
		})
	})
}

func NewHTTPServer(options ...func(*HTTPServer) error) (*HTTPServer, error) {
	store, err := CreateSQLStore("postgres",
		"dbname=blog user=postgres password=postgres host=localhost sslmode=disable",
		squirrel.Dollar)
	if err != nil {
		return nil, err
	}

	uc, err := NewArticleUseCase(store)

	if err != nil {
		return nil, err
	}

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"https://*", "http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	httpServer := &HTTPServer{
		Host:   "127.0.0.1",
		Port:   8000,
		router: r,
		uc:     uc,
	}

	httpServer.setupRoute()

	if len(options) == 0 {
		return httpServer, nil
	}

	for _, opt := range options {
		if err := opt(httpServer); err != nil {
			return nil, err
		}
	}

	return httpServer, nil
}

func (s *HTTPServer) Start() {
	listen := fmt.Sprintf("%s:%d", s.Host, s.Port)

	http.ListenAndServe(listen, s.router)
}

func wrapError(err error) []byte {
	wrapper := struct {
		Message string `json:"message"`
	}{Message: err.Error()}

	j, _ := json.Marshal(wrapper)

	return j
}

func writeError(w http.ResponseWriter, status int, err error) {
	w.Header().Add("content-type", "application-json")
	w.WriteHeader(status)
	w.Write(wrapError(err))
}

func writeJSON(w http.ResponseWriter, status int, result interface{}) {
	w.Header().Add("content-type", "application-json")
	w.WriteHeader(status)

	json.NewEncoder(w).Encode(result)
}

var (
	ErrInvalidRequestPayload = errors.New("invalid request payload")
)

func (s *HTTPServer) NewArticleHandler(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	}
	ctx := r.Context()

	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		writeError(w, http.StatusBadRequest, ErrInvalidRequestPayload)
		return
	}

	article, err := s.uc.CreateArticle(ctx, payload.Title, payload.Content)
	if err != nil {
		writeError(w, http.StatusUnprocessableEntity, err)
		return
	}

	result := struct {
		ID        string `json:"id"`
		CreatedAt string `json:"created_at"`
	}{article.ID.String(), article.CreatedAt.Format(time.RFC3339)}

	writeJSON(w, http.StatusCreated, result)

}

func (s *HTTPServer) SingleArticleHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	aidStr := chi.URLParamFromCtx(ctx, "articleID")

	aid, err := uuid.Parse(aidStr)

	if err != nil {
		writeError(w, http.StatusNotFound, err)
		return
	}

	article, err := s.uc.FindArticleByID(ctx, aid)

	if err != nil {
		if err == ErrArticleNotFound {
			writeError(w, http.StatusNotFound, err)
		} else {
			writeError(w, http.StatusInternalServerError, err)
		}
		return
	}

	writeJSON(w, http.StatusOK, article)
}

func (s *HTTPServer) ListArticlesHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	briefs, err := s.articleReader.ListArticles(ctx)

	if err != nil {
		if err == ErrArticleNotFound {
			writeError(w, http.StatusNotFound, err)
		} else {
			writeError(w, http.StatusInternalServerError, err)
		}
		return
	}

	writeJSON(w, http.StatusOK, briefs)
}

func (s *HTTPServer) DeleteArticleHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	aidStr := chi.URLParamFromCtx(ctx, "articleID")

	aid, err := uuid.Parse(aidStr)

	if err != nil {
		writeError(w, http.StatusNotFound, err)
		return
	}

	err = s.uc.DeleteArticle(ctx, aid)

	if err != nil {
		if err == ErrArticleNotFound {
			writeError(w, http.StatusNotFound, err)
		} else {
			writeError(w, http.StatusInternalServerError, err)
		}
		return
	}

	payload := struct {
		ID string `json:"id"`
	}{
		ID: aid.String(),
	}
	writeJSON(w, http.StatusGone, payload)

}

func (s *HTTPServer) EditArticleHandler(w http.ResponseWriter, r *http.Request) {
	writeError(w, http.StatusNotImplemented, ErrNotImplemented)
}
