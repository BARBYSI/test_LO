package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	rest "task-manager/internal"
	"task-manager/internal/config"
	"task-manager/pkg/logger"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(fmt.Sprintf("failed to load config: %v", err))
	}

	logger.Init(&logger.LoggerConfig{Enabled: cfg.LoggerEnabled})

	rest := rest.NewRest(cfg)
	logger.LogInfo("starting server...")

	shutdownCh := make(chan struct{})
	go func() {
		if err := rest.RunRest(); err != nil {
			logger.LogError(fmt.Sprintf("failed to start server: %v", err))
			close(shutdownCh)
		}
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	select {
	case sig := <-sigCh:
		logger.LogInfo(fmt.Sprintf("received signal %s, shutting down server", sig))
	case <-shutdownCh:
		logger.LogInfo("server shutdown")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := rest.ShutdownRest(ctx); err != nil {
		logger.LogError(fmt.Sprintf("graceful shutdown failed: %v", err))
	} else {
		logger.LogInfo("graceful shutdown succeed")
	}
}
