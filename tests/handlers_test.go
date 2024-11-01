package tests

import (
	"context"
	"os"
	"testing"

	"github.com/FischukSergey/gophkeeper/cmd/server/initial"
	"github.com/FischukSergey/gophkeeper/internal/app/services"
	"github.com/FischukSergey/gophkeeper/tests/suite"
)

var service *services.GRPCService

// TestMain инициализация тестов.
func TestMain(m *testing.M) {
	os.Setenv("CONFIG_PATH", "../config/local.yml")
	initial.InitConfig()
	logger := initial.InitLogger()
	storage := suite.InitTestStorage()
	logger.Info("TestMain database connected")
	service = services.NewGRPCService(logger, storage, nil)
	logger.Info("TestMain service created")
	os.Exit(m.Run())
}

func TestPing(t *testing.T) {
	err := service.Ping(context.Background())
	if err != nil {
		t.Errorf("failed to ping: %v", err)
	}
}
