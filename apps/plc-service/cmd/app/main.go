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
	Pressure float64
	Temp     float64
	Fuel     float64
	Setpoint float64
}

func main() {
	log.Println("PLC Service Starting...")
	client, err := physics.NewClient("localhost:50051")
	if err != nil {
		log.Fatalf("Failed to connect to Physics: %v", err)
	}
	defer client.Close()
	pid := domain.NewPID(5.0, 0.2, 1.0)
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
		setpoint := 60.0
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				status, err := client.GetStatus(ctx)
				if err != nil {
					log.Printf("Error getting status: %v", err)
					continue
				}
				fuelCmd := pid.Update(setpoint, status.SteamPressure, 1.0)
				steamCmd := 0.0
				if status.Timestamp > 10.0 {
					steamCmd = 30.0
				}
				if _, err := client.SetControls(ctx, fuelCmd, 50.0, steamCmd); err != nil {
					log.Printf("Error setting controls: %v", err)
				}
				select {
				case telemetryChan <- Telemetry{
					Pressure: status.SteamPressure,
					Temp:     status.FurnaceTemp,
					Fuel:     fuelCmd,
					Setpoint: setpoint,
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
				opcServer.UpdateData(data.Pressure, data.Temp, data.Fuel, data.Setpoint)
			}
		}
	}()
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop
	log.Println("Shutting down PLC Service...")
	opcServer.Stop()
}

// ^\s*\r?\n
// ^[ \t]+
