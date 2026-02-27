package main

import (
	"context"
	"gateway-service/internal/config"
	"gateway-service/internal/domain"
	"gateway-service/internal/mqtt"
	"gateway-service/internal/opcua"
	"gateway-service/internal/worker"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"
)

func main() {
	cfg := config.Load()

	logger, _ := zap.NewProduction()
	defer logger.Sync()
	sugar := logger.Sugar()

	sugar.Infof("IIoT Gateway starting... (Mocks=%v)", cfg.UseMocks)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var poller opcua.IPoller
	if cfg.UseMocks {
		poller = opcua.NewMockPoller()
	} else {
		poller = opcua.NewOpcClient(cfg.OPCUAEndpoint)
	}

	if err := poller.Connect(ctx); err != nil {
		sugar.Fatalf("Poller connect failed: %v", err)
	}

	var publisher mqtt.IPublisher
	if cfg.UseMocks {
		publisher = mqtt.NewMockPublisher()
	} else {
		publisher = mqtt.NewMqttClient(cfg.MQTTBroker, cfg.MQTTClientID)
	}

	if err := publisher.Connect(); err != nil {
		sugar.Fatalf("MQTT connect failed: %v", err)
	}

	pool := worker.NewPool(publisher, cfg.BufferSize, logger)
	pool.Start(ctx, cfg.WorkerCount)

	go func() {
		ticker := time.NewTicker(time.Duration(cfg.PollIntervalMs) * time.Millisecond)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case t := <-ticker.C:
				tags, err := poller.Read(ctx)
				if err != nil {
					sugar.Warnw("Read error", "err", err)
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

	sugar.Info("Shutting down...")
	cancel()
	pool.Stop()
	publisher.Close()
	poller.Close()
	sugar.Info("Gateway stopped cleanly")
}
