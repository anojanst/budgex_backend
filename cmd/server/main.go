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
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	gdb := db.Must(db.Connect(cfg.PostgresDSN()))

	app := api.Build(gdb)

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
