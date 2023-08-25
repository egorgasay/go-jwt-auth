package bootstrap

import (
	"go-jwt-auth/internal/handler"
	"go-jwt-auth/internal/handler/routes"
	"go-jwt-auth/internal/lib"
	"go-jwt-auth/internal/repository"
	"go-jwt-auth/internal/services"
	"go-jwt-auth/internal/storage"
	"go.uber.org/fx"
)

var CommonModules = fx.Options(
	handler.Module,
	routes.Module,
	lib.Module,
	services.Module,
	storage.Module,
	repository.Module,
)
