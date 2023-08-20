package main

import (
	"github.com/gin-gonic/gin"
	"go-jwt-auth/config"
	"go-jwt-auth/internal/handler"
	"go-jwt-auth/internal/handler/routes"
	"go-jwt-auth/internal/service"
	"go-jwt-auth/internal/storage"
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

	conf, err := config.New()
	if err != nil {
		log.Fatalf("cannot create config instance %v", err)
	}

	st, err := storage.New(conf.StorageConfig, logger)
	if err != nil {
		log.Fatalf("cannot create storage instance %v", err)
	}

	useCase, err := service.NewTokenManager(st, logger, conf.JWTConfig)
	if err != nil {
		log.Fatalf("cannot create service instance %v", err)
	}

	h := handler.New(useCase, logger)
	router := gin.Default()
	routes.Routes(router.Group("/"), h)

	done := make(chan os.Signal)
	signal.Notify(done, os.Interrupt)

	go func() {
		if conf.HTTPS {
			log.Fatal("not implemented") // TODO
		} else {
			if err := http.ListenAndServe(":"+conf.Port, router); err != nil {
				log.Printf("cannot start http server %v", err)
			}
		}
	}()

	<-done

	if err := st.Close(); err != nil {
		log.Printf("cannot close storage %v", err)
	}
}
