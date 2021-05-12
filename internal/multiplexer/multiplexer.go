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

func Run() {
	cfg := config.NewConfig()
	srv := http.NewServer(cfg, handlers.Init(cfg))
	go func() {
		if err := srv.Run(); err != nil {
			log.Errorf("Error occurred while running http server: %v", err.Error())
		} else {
			log.Info("ok")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit

	const timeout = 5 * time.Second

	ctx, shutdown := context.WithTimeout(context.Background(), timeout)
	defer shutdown()

	if err := srv.Stop(ctx); err != nil {
		log.Printf("Failed to stop server: %v", err)
	} else {
		log.Println("Graceful shutdown")
	}
}
