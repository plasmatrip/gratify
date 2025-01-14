package api

import (
	"github.com/plasmatrip/gratify/internal/config"
	"github.com/plasmatrip/gratify/internal/controller"
	"github.com/plasmatrip/gratify/internal/logger"
	"github.com/plasmatrip/gratify/internal/repository"
)

type Dependencies struct {
	Config     config.Config
	Logger     logger.Logger
	Repo       repository.Repository
	Controller controller.Controler
}
