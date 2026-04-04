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

	handlerScenarios "github.com/agrrh/mycorp/internal/handlers/scenarios"
	"github.com/agrrh/mycorp/internal/scenario/store"
)

func main() {
	scenarioDir := os.Getenv("SCENARIO_DIR")
	if scenarioDir == "" {
		scenarioDir = "./scenarios"
	}

	scStore := store.New(scenarioDir)

	if err := scStore.Load(); err != nil {
		log.Fatal(err)
	}

	scHandler := handlerScenarios.Handler{
		ScStore: scStore,
	}

	// TODO: Add SSO authentication

	e := echo.New()
	e.Use(middleware.RequestLogger())
	// e.Use(middleware.RemoveTrailingSlash())
	e.Use(middleware.Recover())

	// Routes
	e.GET("/scenarios", scHandler.List)
	e.GET("/scenarios/:namespace", scHandler.ListByNamespace)
	e.GET("/scenarios/:namespace/:name/_cli", scHandler.GetCLI)
	e.POST("/scenarios/:namespace/:name", scHandler.Run)

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
