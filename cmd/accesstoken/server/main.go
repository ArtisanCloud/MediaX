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

const (
	defaultListenHost = "127.0.0.1"
	defaultListenPort = "7071"
	defaultListenAddr = defaultListenHost + ":" + defaultListenPort
)

func main() {
	configFlag := flag.String("config", "", "Path to config.yaml (defaults to MEDIA_X_CONFIG or config.yaml)")
	listenAddrFlag := flag.String("listen-addr", "", "Full listen address host:port (overrides ACCESSTOKEN_LISTEN_ADDR)")
	portFlag := flag.String("port", "", "Listen port or address, e.g. 7071 or :7071 (overrides only the port component)")
	flag.Parse()

	if err := run(*configFlag, *listenAddrFlag, *portFlag); err != nil {
		log.Fatalf("accesstoken-server: %v", err)
	}
}

func run(configFlag, listenAddrFlag, portFlag string) error {
	defaultConfigPath := app.ResolveConfigPath(configFlag)

	listenAddr := resolveListenAddr(listenAddrFlag, portFlag)

	server, cleanup, err := newAccessTokenServer(defaultConfigPath, listenAddr)
	if err != nil {
		return err
	}
	defer cleanup()
	if err := server.validateListenAddr(); err != nil {
		return err
	}

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

	server.logServerEvent("listening", map[string]any{
		"api_token":          app.MaskToken(server.apiToken),
		"config":             defaultConfigPath,
		"flow_ttl_seconds":   server.flowTTLSeconds,
		"redis_available":    server.storageBackend == storageBackendRedis,
		"warning_public_net": server.listenAddrPublic,
	})
	if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("listen and serve: %w", err)
	}
	server.logger.InfoF("accesstoken-server: stopped")
	return nil
}
