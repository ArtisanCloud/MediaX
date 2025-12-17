package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	app "github.com/ArtisanCloud/MediaX/cmd/accesstoken/internal/app"
)

const defaultListenAddr = ":7071"

func main() {
	configFlag := flag.String("config", "", "Path to config.yaml (defaults to MEDIA_X_CONFIG or config.yaml)")
	portFlag := flag.String("port", "", "Listen port or address, e.g. 7071 or :7071")
	flag.Parse()

	if err := run(*configFlag, *portFlag); err != nil {
		log.Fatalf("accesstoken-server: %v", err)
	}
}

func run(configFlag, portFlag string) error {
	defaultConfigPath := app.ResolveConfigPath(configFlag)

	listenAddr := resolveListenAddr(portFlag)

	server, cleanup, err := newAccessTokenServer(defaultConfigPath, listenAddr)
	if err != nil {
		return err
	}
	defer cleanup()

	httpServer := &http.Server{
		Addr:         listenAddr,
		Handler:      server.routes(),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 60 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()
		if err := httpServer.Shutdown(shutdownCtx); err != nil && !errors.Is(err, http.ErrServerClosed) {
			server.logger.ErrorF("accesstoken-server: graceful shutdown failed: %v", err)
		}
	}()

	server.logger.InfoF("accesstoken-server: listening addr=%s api_token=%s config=%s", listenAddr, app.MaskToken(server.apiToken), defaultConfigPath)
	if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("listen and serve: %w", err)
	}
	server.logger.InfoF("accesstoken-server: stopped")
	return nil
}
