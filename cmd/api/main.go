package main

import (
	"flag"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"log"
	"snakey/pkg/server"
	"snakey/pkg/service"
	"syscall"
	"time"
)

var (
	addr            = flag.String("addr", ":8080", "Application http server network address")
	shutdownTimeout = flag.Duration("shutdown-timeout", 30*time.Second, "Graceful shutdown timeout")
)

func main() {
	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	api := NewSnakeyAPI()

	r.Route("/", func(r chi.Router) {
		r.Post("/new", api.New)
		r.Post("/validate", api.Validate)
		r.Get("/health", api.Health)
	})

	s := server.NewHTTPServer(*addr, r, *shutdownTimeout)

	s.Start()
	log.Printf("started api server at %s", *addr)
	sig := service.Wait(syscall.SIGTERM, syscall.SIGINT)
	log.Printf("recieved signal %s", sig.String())
	s.Stop()
}
