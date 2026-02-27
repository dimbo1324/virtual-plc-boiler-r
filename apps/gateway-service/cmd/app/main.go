package main

import (
	"context"
	"gateway-service/internal/domain"
	"gateway-service/internal/mqtt"
	"gateway-service/internal/opcua"
	"gateway-service/internal/worker"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	USE_MOCKS    = true
	WORKER_COUNT = 5
	BUFFER_SIZE  = 100
)

func main() {
	log.Printf("IIoT Gateway (Mocks=%v) starting...", USE_MOCKS)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var poller opcua.IPoller
	var publisher mqtt.IPublisher

	if USE_MOCKS {
		poller = &opcua.MockPoller{}
		publisher = &mqtt.MockPublisher{}
	} else {
		poller = opcua.NewOpcClient("opc.tcp://localhost:4840")
		publisher = mqtt.NewMqttClient("tcp://localhost:1883", "gw_01")
	}

	if err := poller.Connect(ctx); err != nil {
		log.Fatalf("Poller error: %v", err)
	}
	if err := publisher.Connect(); err != nil {
		log.Fatalf("Publisher error: %v", err)
	}

	pool := worker.NewPool(publisher, BUFFER_SIZE)
	pool.Start(ctx, WORKER_COUNT)

	go func() {
		ticker := time.NewTicker(500 * time.Millisecond)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case t := <-ticker.C:
				tags, err := poller.Read(ctx)
				if err != nil {
					log.Printf("Read error: %v", err)
					continue
				}

				payload := domain.Payload{
					Timestamp: t.UTC().Format(time.RFC3339),
					AssetID:   "boiler_01",
					Tags:      tags,
				}

				pool.Push(payload)
			}
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	log.Println("Shutting down...")
	cancel()
	pool.Stop()
	publisher.Close()
	poller.Close()
}
