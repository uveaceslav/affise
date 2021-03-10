package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Server struct {
	logger *log.Logger
	router *http.ServeMux
}

func NewServer(
	logger *log.Logger,
	router *http.ServeMux,
) *Server {
	return &Server{
		logger: logger,
		router: router,
	}
}

func (s *Server) Serve(address string) {
	server := &http.Server{
		Addr:    address,
		Handler: s.router,
	}

	done := make(chan bool)
	quit := make(chan os.Signal, 1)
	// signal.Notify(quit, os.Interrupt)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-quit
		s.logger.Println("Server is shutting down...")

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		server.SetKeepAlivesEnabled(false)
		if err := server.Shutdown(ctx); err != nil {
			s.logger.Fatalf("Could not gracefully shutdown the server: %v\n", err)
		}
		close(done)
	}()

	s.logger.Println("Server is ready to handle requests at", address)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		s.logger.Fatalf("Could not listen on %s: %v\n", address, err)
	}

	<-done
	s.logger.Println("Server stopped")
}
