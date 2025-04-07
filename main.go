package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"log"
	"net/http"
	"os"
	"os/signal"
	"osm-tail/http_handler"
	"osm-tail/router"
	"osm-tail/utils/envconf"
	"osm-tail/utils/postgresql"
	"syscall"
	"time"
)

func main() {
	fmt.Println("Hello, World!")

	err := envconf.LoadAppConfig()

	if err != nil {
		log.Fatalf("Error initializing env variables: %v", err)
	}

	db, err := postgresql.NewPostgreSQL(
		envconf.App.PostgreSQL.Host,
		envconf.App.PostgreSQL.Port,
		envconf.App.PostgreSQL.User,
		envconf.App.PostgreSQL.Password,
		envconf.App.PostgreSQL.Database,
		envconf.App.PostgreSQL.LogLevel,
	)

	if err != nil {
		log.Println("Error connecting to database", err)
		return
	}

	r := gin.Default()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	httpHandler := http_handler.Handler{
		Tracer: otel.Tracer("handler"),
		Db:     db,
	}

	router.RegisterRoute(r, httpHandler)

	addr := fmt.Sprintf(":%d", 3000)
	srv := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		<-c
		log.Println("Shutting down server...")

		shutdownCtx, shutdownCancel := context.WithTimeout(ctx, 5*time.Second)
		defer shutdownCancel()

		if err := srv.Shutdown(shutdownCtx); err != nil {
			log.Fatalf("Server forced to shutdown: %v", err)
		}

		log.Println("Server stopped gracefully.")
		cancel()
	}()

	log.Printf("Server start at: http://localhost%s\n", addr)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Failed to start server: %v", err)
	}

	<-ctx.Done()
	log.Println("Program exited gracefully.")
}
