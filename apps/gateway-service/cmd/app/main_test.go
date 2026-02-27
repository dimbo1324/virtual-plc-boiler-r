package main_test

import (
	"gateway-service/internal/config"
	"os"
	"syscall"
	"testing"
	"time"
)

func TestMainIntegration(t *testing.T) {
	oldLoad := config.Load
	config.Load = func() *config.Config {
		return &config.Config{
			UseMocks:       true,
			WorkerCount:    1,
			BufferSize:     10,
			PollIntervalMs: 100,
		}
	}
	defer func() { config.Load = oldLoad }()

	go main_main()

	time.Sleep(500 * time.Millisecond)
	os.Signal(syscall.SIGTERM)
}
