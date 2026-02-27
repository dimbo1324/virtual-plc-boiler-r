package main

import (
	"os"
	"syscall"
	"testing"
	"time"

	"gateway-service/internal/config"
)

func TestMainIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("integration test skipped in short mode")
	}

	originalLoad := config.Load
	defer func() { config.Load = originalLoad }()

	config.Load = func() *config.Config {
		return &config.Config{
			UseMocks:       true,
			WorkerCount:    2,
			BufferSize:     20,
			PollIntervalMs: 100,
		}
	}

	done := make(chan struct{})

	go func() {
		defer close(done)
		main()
	}()

	time.Sleep(800 * time.Millisecond)

	p, err := os.FindProcess(os.Getpid())
	if err != nil {
		t.Fatal(err)
	}
	if err := p.Signal(syscall.SIGTERM); err != nil {
		t.Fatal(err)
	}

	select {
	case <-done:
		t.Log("Main exited cleanly")
	case <-time.After(4 * time.Second):
		t.Fatal("Timeout: main() не остановился после SIGTERM")
	}
}
