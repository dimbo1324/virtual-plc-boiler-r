package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"plc-service/internal/clients/physics"
	"plc-service/internal/domain"
)

const (
	PhysicsAddress = "localhost:50051"
	TargetPressure = 60.0
	TargetLevel    = 500.0
)

func main() {
	log.Println("PLC Service Starting...")

	physClient, err := physics.NewClient(PhysicsAddress)
	if err != nil {
		log.Fatalf("Failed to connect to physics: %v", err)
	}
	defer physClient.Close()
	log.Println("Connected to Physics Engine")

	pressurePID := domain.NewPID(5.0, 0.1, 10.0)

	ctx, cancel := context.WithCancel(context.Background())

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		ticker := time.NewTicker(100 * time.Millisecond)
		defer ticker.Stop()

		var lastTime = time.Now()

		for {
			select {
			case <-ctx.Done():
				return
			case t := <-ticker.C:
				dt := t.Sub(lastTime).Seconds()
				lastTime = t

				status, err := physClient.GetStatus(ctx)
				if err != nil {
					log.Printf("Error reading sensors: %v", err)
					continue
				}

				fuelCmd := pressurePID.Update(TargetPressure, status.SteamPressure, dt)

				levelError := TargetLevel - status.DrumLevel
				waterCmd := 50.0 + (levelError * 0.5)

				if waterCmd > 100 {
					waterCmd = 100
				}
				if waterCmd < 0 {
					waterCmd = 0
				}

				steamCmd := 0.0
				if status.Timestamp > 20.0 {
					steamCmd = 40.0
				}

				_, err = physClient.SetControls(ctx, fuelCmd, waterCmd, steamCmd)
				if err != nil {
					log.Printf("Error writing controls: %v", err)
				}

				if t.UnixMilli()%1000 < 100 {
					fmt.Printf("[PLC] T: %.1fs | Press: %.2f Bar (Set: %.0f) | Fuel: %.1f%% | Level: %.0f mm | Flow: %.1f\n",
						status.Timestamp, status.SteamPressure, TargetPressure, fuelCmd, status.DrumLevel, status.SteamFlow)
				}
			}
		}
	}()

	<-sigChan
	log.Println("\nShutting down PLC...")
	cancel()
	time.Sleep(500 * time.Millisecond)
	log.Println("Bye!")
}
