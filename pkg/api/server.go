package api

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

type Server struct {
	httpServer *http.Server
	router     *http.ServeMux
	host       string
	port       string
	quit       chan bool
}

func NewServer(port string, dbPath string, dbSchemaPath string, geoSeedPath string) (*Server, error) {
	router := http.NewServeMux()

	httpServer := &http.Server{Addr: port, Handler: router}

	// router.HandleFunc("/detect", endpoints.WithLogging(dh.Detect))
	return &Server{
		router:     router,
		port:       port,
		httpServer: httpServer,
		quit:       make(chan bool),
	}, nil
}

func (s *Server) Start() error {
	log.Printf("Starting http service on %s", s.port)
	// Start server
	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := s.Stop(ctx); err != nil {
		return err
	}
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	log.Println("gracefully shutting down")

	return s.httpServer.Shutdown(ctx)
}
