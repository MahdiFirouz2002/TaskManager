package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	concurrentlogger "nikandishan/concurrentLogger"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

func gracefulShutdown(srv *http.Server, wg *sync.WaitGroup) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	fmt.Println("shout down signal received. shoting down the server...")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Error shutting down server: %v", err)
	}

	close(concurrentlogger.UpdateTaskChan)
	wg.Wait()
	fmt.Println("server gracefully stoped")
}

func StartServer() *http.Server {
	var wg sync.WaitGroup
	concurrentlogger.StartLogger(&wg)

	r := gin.Default()

	r.Use(LoggingMiddleware())

	setupRoutes(r)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	go func() {
		log.Println("Server started on port 8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	gracefulShutdown(srv, &wg)
	return srv
}
