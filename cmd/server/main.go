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

	"github.com/agrrh/mycorp/internal/application/server/config"
	"github.com/agrrh/mycorp/internal/application/server/handlers"
	"github.com/agrrh/mycorp/internal/domain/scenario_store"
)

func main() {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "./examples/server.config.yaml"
	}

	cfg, err := config.Load(configPath)
	if err != nil {
		log.Fatal(err)
	}

	scenarioDir := os.Getenv("SCENARIO_DIR")
	if scenarioDir == "" {
		scenarioDir = "./examples/scenarios"
	}

	scStore := scenario_store.New(scenarioDir)

	if err := scStore.Load(); err != nil {
		log.Fatal(err)
	}

	sHandler := handlers.Handler{
		ScStore: scStore,
	}

	e := echo.New()
	e.Use(middleware.RequestLogger())
	// e.Use(middleware.RemoveTrailingSlash())
	e.Use(middleware.Recover())

	// Routes
	scenarios := e.Group("/scenarios")

	// TODO: Use SSO authentication
	scenarios.GET("/", sHandler.List)
	scenarios.GET("/:namespace", sHandler.ListByNamespace)
	scenarios.GET("/:namespace/:name/_cli", sHandler.GetCLI)
	scenarios.POST("/:namespace/:name", sHandler.Run)

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
