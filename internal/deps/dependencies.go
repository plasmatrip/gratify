package deps

import (
	"github.com/plasmatrip/gratify/internal/api/middleware/logger"
	"github.com/plasmatrip/gratify/internal/config"

	"github.com/plasmatrip/gratify/internal/repository"
)

type Dependencies struct {
	Config config.Config
	Logger logger.Logger
	Repo   repository.Repository
}
