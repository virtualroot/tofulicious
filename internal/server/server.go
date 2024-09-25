package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"golang.org/x/sync/errgroup"
)

type HTTPServer struct {
	server *http.Server
}

func NewHTTPServer() (*HTTPServer, error) {
	s := &HTTPServer{}
	mux := chi.NewRouter()
	s.server = &http.Server{
		Addr:              fmt.Sprintf("%s:%d", "localhost", 8080),
		Handler:           mux,
		ReadHeaderTimeout: time.Minute,
	}
	mux.Use(middleware.RequestID)
	mux.Use(middleware.Logger)
	mux.Use(middleware.Recoverer)
	mux.Use(render.SetContentType(render.ContentTypeJSON))

	mux.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})

	mux.Get("/api/v1/state/default", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(""))
	})

	// Persist the state to backend
	// ?ID=964143dc-5bf1-1530-34a7-7dc85c5079e7
	mux.Post("/api/v1/state/default", func(w http.ResponseWriter, r *http.Request) {
		ID := r.URL.Query().Get("ID")
		w.Write([]byte("Query: " + ID))

		w.Write([]byte(""))
	})

	// Lock the state
	// When locking support is enabled it will use LOCK and UNLOCK requests providing the lock info in the body.
	// The endpoint should return a 423: Locked or 409: Conflict with the holding lock info when it's already taken, 200: OK for success.
	mux.Put("/api/v1/state/default/lock", func(w http.ResponseWriter, r *http.Request) {
		render.Status(r, http.StatusCreated)
		w.WriteHeader(200)
		w.Write([]byte("root."))
	})

	// Unlock the state
	mux.Delete("/api/v1/state/default/lock", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte(""))
	})

	return s, nil
}

// Start start the HTTP server.
func (s *HTTPServer) Start() error {
	errg, _ := errgroup.WithContext(context.Background())
	errg.Go(func() error {
		log.Print("Starting server ", s.server.Addr)
		err := s.server.ListenAndServe()
		if err != http.ErrServerClosed {
			return err
		}
		return nil

	})
	return errg.Wait()
}

// Shutdown gracefully shut down the HTTP server.
func (s *HTTPServer) Shutdown(ctx context.Context) error {
	log.Print("Stopping server ", s.server.Addr)
	return s.server.Shutdown(ctx)
}
