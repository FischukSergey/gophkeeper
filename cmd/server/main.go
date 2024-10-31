package main

import (
	"fmt"

	"github.com/FischukSergey/gophkeeper/cmd/server/initial"
	"github.com/FischukSergey/gophkeeper/internal/app/server"
)

var (
	buildVersion = "N/A" // версия сборки
	buildDate    = "N/A" // дата сборки
	buildCommit  = "N/A" // коммит сборки
)

func main() {
	fmt.Printf("Build version: %s\n", buildVersion)
	fmt.Printf("Build date: %s\n", buildDate)
	fmt.Printf("Build commit: %s\n", buildCommit)

	initial.InitConfig() // инициализация конфигурации
	log := initial.InitLogger()
	application := server.NewGrpcServer(log, initial.Cfg.GRPC.Port)
	application.MustRun()
}
