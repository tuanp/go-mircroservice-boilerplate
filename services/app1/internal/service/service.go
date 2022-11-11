package service

import (
	"github.com/rogpeppe/go-internal/cache"
	"github.com/tuanp/go-mircroservice-boilerplate/pkg/logger"
	"github.com/tuanp/go-mircroservice-boilerplate/services/app1/internal/repository"
)

type Services struct {
}

type Deps struct {
	Repos       *repository.Repositories
	Cache       cache.Cache
	CacheTTL    int64
	Environment string
	Domain      string
	Logger      logger.Logger
}

func NewServices(deps Deps) *Services {
	// Services init
	//schoolsService := NewSchoolsService(deps.Repos.Schools, deps.Cache, deps.CacheTTL)

	return &Services{
		//Schools:        schoolsService,

	}
}