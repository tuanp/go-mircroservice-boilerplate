package app

import (
	"log"

	"github.com/tuanp/go-gin-boilerplate/config"
)

// Run initializes whole application.
func Run(configPath string) {
	log.Println("Starting api server")
	cfg, err := config.Init(configPath)
	//if err != nil {
	//	logger.Error(err)

	//	return
	//}
}
