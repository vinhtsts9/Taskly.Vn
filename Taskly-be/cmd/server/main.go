package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"Taskly.com/m/internal/initialize"
	"Taskly.com/m/internal/middlewares"
	websocket "Taskly.com/m/ws"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
)

var interruptSignals = []os.Signal{
	os.Interrupt,
	syscall.SIGTERM,
	syscall.SIGINT,
}

func main() {
	r := initialize.Run()

	cm := websocket.NewConnectionManager()

	r.GET("/v1/2024/ws", middlewares.AuthenMiddleware(), middlewares.CasbinMiddleware(), func(c *gin.Context) {
		websocket.HandleConnections(c, cm)
	})

	r.GET("/checkStatus", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	ctx, stop := signal.NotifyContext(context.Background(), interruptSignals...)
	defer stop()

	waitGroup, ctx := errgroup.WithContext(ctx)
	go cm.Run(ctx)
	runGinServer(ctx, waitGroup, r)

	err := waitGroup.Wait()
	if err != nil {
		log.Fatalf("Error from wait group: %v", err)
	}

	log.Println("Server has stopped")
}

func runGinServer(ctx context.Context, waitGroup *errgroup.Group, r *gin.Engine) {
	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	waitGroup.Go(func() error {
		log.Println("Starting server on :8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
			return err
		}
		return nil
	})

	waitGroup.Go(func() error {
		<-ctx.Done()
		log.Println("Graceful shutdown HTTP server")
		if err := srv.Shutdown(context.Background()); err != nil {
			log.Printf("HTTP server shutdown error: %v", err)
		}
		return nil
	})
}
