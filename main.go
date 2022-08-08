package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/awcjack/cloudbet/application"
	"github.com/awcjack/cloudbet/infrastructure"
	"github.com/awcjack/cloudbet/interfaces"
	"github.com/sirupsen/logrus"
)

func main() {
	repo := infrastructure.NewMemoryRepository()
	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)

	app := application.NewApplication(repo, repo, repo, repo, logger)

	httpServer := interfaces.NewHttpServer(*app)

	server := &http.Server{
		Addr:    ":8080",
		Handler: interfaces.NewHandler(*httpServer),
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
			log.Printf("listen: %s\n", err)
		}
	}()

	cloudbetCrawler := application.NewCloudbetHander(repo, logger, "***REMOVED***")

	ticker := time.NewTicker(5 * time.Second)
	exitTicker := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				cloudbetCrawler.StoreAllEvents(context.Background())
				cloudbetCrawler.CheckEventsCloseToCutOff(context.Background())
				// do stuff
			case <-exitTicker:
				ticker.Stop()
				return
			}
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	close(exitTicker)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}
}
