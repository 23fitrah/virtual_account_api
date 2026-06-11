package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"virtual_account_api/internal/injector"
	"virtual_account_api/internal/routes"
	"virtual_account_api/utils"
	"syscall"
	"time"
)

func main() {
	// Init Wire
	container, err := injector.InitializeApp()
	if err != nil {
		utils.LogError("[FAILED] Failed to initialize app", err)
	}

	// Setup Router
	router := routes.SetupRouter(container)

	// Running Graceful Shutdown
	srv := &http.Server{
		Addr:    ":" + os.Getenv("PORT"),
		Handler: router,
	}

	// Running Server at Go Routine
	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			utils.LogError("[FAILED] Failed to start server", err)
		}
	}()

	// Get Signal Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	utils.LogInfo("[INFO] Received termination, shutting down server")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		utils.LogError("[FAILED] Failed to shutdown server", err)
	}

	utils.LogInfo("[INFO] Server exited gracefully")
}
