package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	server "github.com/agrrh/mycorp/internal/application/server"
	"github.com/agrrh/mycorp/internal/domain/scenario_store"
)

func main() {
	scenarioDir := os.Getenv("SCENARIO_DIR")
	if scenarioDir == "" {
		scenarioDir = "./scenarios"
	}

	scStore := scenario_store.New(scenarioDir)

	if err := scStore.Load(); err != nil {
		log.Fatal(err)
	}

	sHandler := server.Handler{
		ScStore: scStore,
	}

	// TODO: Add SSO authentication

	e := echo.New()
	e.Use(middleware.RequestLogger())
	// e.Use(middleware.RemoveTrailingSlash())
	e.Use(middleware.Recover())

	// Routes
	e.GET("/scenarios", sHandler.List)
	e.GET("/scenarios/:namespace", sHandler.ListByNamespace)
	e.GET("/scenarios/:namespace/:name/_cli", sHandler.GetCLI)
	e.POST("/scenarios/:namespace/:name", sHandler.Run)

	// Start server
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: e,
	}
	go func() {
		if err := e.StartServer(srv); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()
	log.Printf("server started on :%s", port)

	<-ctx.Done()
	log.Println("shutting down...")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(shutdownCtx); err != nil {
		log.Fatal(err)
	}
}
