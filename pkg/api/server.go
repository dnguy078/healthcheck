package api

import (
	"context"
	"log"
	"net/http"

	"github.com/dnguy078/healthcheck/pkg/storage"
)

// Server is a http server
type Server struct {
	httpServer *http.Server
	router     *http.ServeMux
	addr       string
	sslCert    string
	sslKey     string
	quit       chan bool
}

// NewServer returns a http server
func NewServer(addr string, sslCert string, sslKey string, db *storage.Collection) (*Server, error) {
	router := http.NewServeMux()
	httpServer := &http.Server{Addr: addr, Handler: router}

	hh := &HealthCheckHandler{db}

	router.Handle("/api/health/checks/", hh)
	router.Handle("/api/health/checks", hh)

	return &Server{
		router:     router,
		addr:       addr,
		sslCert:    sslCert,
		sslKey:     sslKey,
		httpServer: httpServer,
	}, nil
}

// Start starts the server
func (s *Server) Start() {
	log.Printf("Starting http service on %s", s.addr)
	go func() {
		if err := s.httpServer.ListenAndServeTLS(s.sslCert, s.sslKey); err != nil {
			// if err := s.httpServer.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()
}

// Stop gracefully stops the server
func (s *Server) Stop(ctx context.Context) error {
	log.Println("Stopping http server")
	return s.httpServer.Shutdown(ctx)
}
