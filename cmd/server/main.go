// @title           Budgex API
// @version         0.1.0
// @description     Backend API for Budgex (transactions, categories, budgets).
// @BasePath        /api
// @schemes         http
// @securityDefinitions.apikey BearerAuth
// @in              header
// @name            Authorization
// @description     "Format: Bearer {token}"

package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"budgex_backend/internal/api"
	"budgex_backend/internal/config"
	"budgex_backend/internal/db"
	_ "budgex_backend/internal/docs" // generated package
	"budgex_backend/internal/observability"

	"github.com/clerk/clerk-sdk-go/v2"
	fiberSwagger "github.com/swaggo/fiber-swagger"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	// Logger
	_, err = observability.InitLogger(
		getEnv("SERVICE_NAME", "budgex-backend"),
		getEnv("LOG_LEVEL", "info"),
	)
	if err != nil {
		log.Fatalf("logger: %v", err)
	}

	// Tracing
	tracerCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	tp, err := observability.InitTracer(tracerCtx, getEnv("SERVICE_NAME", "budgex-backend"))
	if err != nil {
		log.Fatalf("tracer: %v", err)
	}
	defer observability.ShutdownTracer(context.Background(), tp)

	// Initialize Clerk SDK
	clerkSecretKey := os.Getenv("CLERK_SECRET_KEY")
	if clerkSecretKey == "" {
		log.Fatal("CLERK_SECRET_KEY environment variable is not set")
	}
	log.Printf("Clerk SDK initialized with secret key: %s...", clerkSecretKey[:10])
	clerk.SetKey(clerkSecretKey)

	gdb := db.Must(db.Connect(cfg.PostgresDSN()))
	if err := db.AutoMigrate(gdb); err != nil {
		log.Fatalf("migrate: %v", err)
	}

	app := api.Build(gdb)

	// Swagger UI (served at /swagger/index.html)
	app.Get("/swagger/*", fiberSwagger.WrapHandler)

	// run server in a goroutine for graceful shutdown
	addr := fmt.Sprintf(":%d", cfg.Port)
	go func() {
		if err := app.Listen(addr); err != nil {
			log.Printf("server stopped: %v", err)
		}
	}()
	log.Printf("server listening on %s", addr)

	// wait for SIGINT/SIGTERM
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	// graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := app.Shutdown(); err != nil {
		log.Printf("shutdown error: %v", err)
	}
	// close DB
	if sqlDB, err := gdb.DB(); err == nil {
		_ = sqlDB.Close()
	}

	<-ctx.Done()
	log.Println("server exited")
}

func getEnv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}
