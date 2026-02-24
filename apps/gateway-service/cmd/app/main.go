package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gateway-service/internal/domain"
	"gateway-service/internal/mqtt"
	"gateway-service/internal/opcua"
	"gateway-service/internal/worker"
)

func main() {
	log.Println("🌐 IIoT Gateway Starting...")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	poller := opcua.NewOpcClient("opc.tcp://localhost:4840")
	publisher := mqtt.NewMqttClient("tcp://localhost:1883", "gateway_01")

	pool := worker.NewPool(publisher, 1000)

	pool.Start(ctx, 3)

	go func() {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case t := <-ticker.C:

				payload := domain.Payload{
					Timestamp: t.UTC().Format(time.RFC3339),
					AssetID:   "boiler_01",
					Tags: domain.Tags{
						Temperature: 450.5,
						Pressure:    60.2,
					},
				}

				pool.Push(payload)
			}
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	log.Println("Shutting down Gateway...")
	cancel()
	pool.Stop()
	log.Println("Bye!")
}
