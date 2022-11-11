package main

import (
	"github.com/tuanp/go-mircroservice-boilerplate/services/app1/internal/app"
)

const configsDir = "config"

func main() {
	app.Run(configsDir)
}
