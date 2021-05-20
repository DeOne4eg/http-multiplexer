package multiplexer

import (
	"context"
	"github.com/DeOne4eg/http-multiplexer/config"
	"github.com/DeOne4eg/http-multiplexer/internal/channel/http"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// Run launches the application.
func Run() {
	// create config
	cfg := config.NewConfig()
	log.Printf("Config: %+v", *cfg)

	// create handlers
	handlers := http.NewHandler()

	// create HTTP server
	srv := http.NewServer(cfg, handlers.Init())

	// run HTTP server
	go func() {
		log.Printf("Running HTTP server on :%d", cfg.HTTP.Port)
		_ = srv.Run()
	}()

	// catch signals for quit from application
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// graceful shutdown with limit 5 seconds
	ctx, shutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdown()

	// if server stopped within 5 seconds then ok else throw error
	if err := srv.Stop(ctx); err != nil {
		log.Printf("Failed to stop server: %v", err)
	} else {
		log.Println("Graceful shutdown")
	}
}
