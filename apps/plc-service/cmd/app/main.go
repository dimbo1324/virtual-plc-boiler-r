package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"plc-service/internal/clients/physics"
	"plc-service/internal/domain"
	"plc-service/internal/opcua"
	"syscall"
	"time"
)

type Telemetry struct {
	Pressure  float64
	Temp      float64
	Fuel      float64
	DrumLevel float64
	SteamFlow float64
}

func main() {
	log.Println("PLC Service Starting...")
	physicsAddr := os.Getenv("PHYSICS_ADDR")
	if physicsAddr == "" {
		physicsAddr = "localhost:50051"
	}
	client, err := physics.NewClient(physicsAddr)
	if err != nil {
		log.Fatalf("Failed to connect to Physics: %v", err)
	}
	defer client.Close()

	pressurePid := domain.NewPID(5.0, 0.2, 1.0)
	levelPid := domain.NewPID(2.0, 0.1, 0.5)

	opcServer := opcua.NewServer(4840)
	telemetryChan := make(chan Telemetry, 1)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		if err := opcServer.Start(ctx); err != nil {
			log.Printf("OPC UA Error: %v", err)
		}
	}()

	go func() {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()

		lastTick := time.Now()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				now := time.Now()
				dt := now.Sub(lastTick).Seconds()
				if dt <= 0 {
					dt = 1.0
				}
				lastTick = now

				status, err := client.GetStatus(ctx)
				if err != nil {
					log.Printf("Error getting status: %v", err)
					continue
				}

				pressureSetpoint := opcServer.GetSetpoint()
				waterSetpoint := 500.0

				fuelCmd := pressurePid.Update(pressureSetpoint, status.SteamPressure, dt)
				waterCmd := levelPid.Update(waterSetpoint, status.DrumLevel, dt)

				startTime := time.Now()

				steamCmd := 0.0
				if time.Since(startTime).Seconds() > 10.0 {
					steamCmd = 30.0
				}
				if status.Timestamp > 10.0 {
					steamCmd = 30.0
				}

				if _, err := client.SetControls(ctx, fuelCmd, waterCmd, steamCmd); err != nil {
					log.Printf("Error setting controls: %v", err)
				}

				select {
				case telemetryChan <- Telemetry{
					Pressure:  status.SteamPressure,
					Temp:      status.FurnaceTemp,
					Fuel:      fuelCmd,
					DrumLevel: status.DrumLevel,
					SteamFlow: status.SteamFlow,
				}:
				default:
				}
			}
		}
	}()

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case data := <-telemetryChan:
				opcServer.UpdateData(data.Pressure, data.Temp, data.Fuel, data.DrumLevel, data.SteamFlow)
			}
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	log.Println("Shutting down PLC Service...")
	opcServer.Stop()
}
