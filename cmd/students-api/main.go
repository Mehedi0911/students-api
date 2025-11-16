package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/mehedi0911/students-api/internal/config"
	"github.com/mehedi0911/students-api/internal/http/handlers/students"
)

func main() {
	//load config
	cfg := config.MustLoad()

	//database setup

	// set up router

	router := http.NewServeMux()
	router.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to students api"))
	})

	router.HandleFunc("POST /api/students", students.New())

	// set up server

	server := http.Server{
		Addr:    cfg.Addr,
		Handler: router,
	}
	slog.Info("Server Started at", slog.String("address", cfg.Addr))

	done := make(chan os.Signal, 1)

	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			log.Fatal("Failed to start server!")
		}
	}()

	<-done

	slog.Info("Shuting down the server!!")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		slog.Error("Failed to shutdown the server!!", slog.String("error", err.Error()))
	}

	slog.Info("Server shutdown successfully!")

}
