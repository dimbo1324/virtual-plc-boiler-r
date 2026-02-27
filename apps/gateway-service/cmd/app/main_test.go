package main

import (
	"os"
	"runtime"
	"testing"
	"time"
)

func TestMainIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("integration test skipped in short mode")
	}
	if runtime.GOOS == "windows" {
		t.Skip("TestMainIntegration: SIGTERM не поддерживается на Windows (используй Linux/Mac для полного теста)")
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
	if err := p.Signal(os.Interrupt); err != nil {
		t.Fatal(err)
	}

	select {
	case <-done:
		t.Log("Main завершился корректно (graceful shutdown)")
	case <-time.After(4 * time.Second):
		t.Fatal("Таймаут: main() не остановился после сигнала")
	}
}
