package main

import (
	"github.com/gin-gonic/gin"
	"go-jwt-auth/internal/handler"
	"go-jwt-auth/internal/storage"
	"go-jwt-auth/internal/usecase"
	"go.uber.org/zap"
	"log"
	"net/http"
	"os"
	"os/signal"
)

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("cannot create logger instance %v", err)
	}

	storage, err := storage.New(logger)
	if err != nil {
		log.Fatalf("cannot create storage instance %v", err)
	}

	useCase := usecase.New(storage, logger)

	handler := handler.New(useCase, logger)
	router := gin.Default()

	done := make(chan os.Signal)
	signal.Notify(done, os.Interrupt)

	go func() {
		// up http server
		if err := http.ListenAndServe(":8080", router); err != nil {
			log.Printf("cannot start http server %v", err)
		}
	}()

	<-done

	if err := storage.Close(); err != nil {
		log.Printf("cannot close storage %v", err)
	}
}
