package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
	"wind/src/core"
	"wind/src/core/interfaces"
)

func main() {
	r := bufio.NewReader(os.Stdin)
	fmt.Print("Wind Digital Twin v1.0 — выберите локацию (1=New York) → ")
	in, _ := r.ReadString('\n')
	in = strings.TrimSpace(in)

	loc := interfaces.Locations[1]
	if in == "1" {
		loc = interfaces.Locations[1]
	} else {
		fmt.Println("Локация не найдена — используется New York")
	}
	runSimulation(loc)
}

func runSimulation(loc interfaces.LocationProfile) {
	location, err := time.LoadLocation(loc.Timezone)
	if err != nil {
		fmt.Printf("Таймзона %s не найдена — использую UTC\n", loc.Timezone)
		location = time.UTC
	}
	now := time.Now().In(location)
	env := core.GenerateWeatherSimulation(loc, now)

	fmt.Printf("ОТЧЁТ: %s, %s | Координаты: %.4f, %.4f | Дата: %d %s | Время: %s (%s)\n",
		loc.City, loc.Country, loc.Latitude, loc.Longitude,
		env.Timestamp.Day(), env.Timestamp.Month().String(), env.Timestamp.Format("15:04:05"), env.TimeOfDay)

	fmt.Printf("Погода: %s | Темп: %.2f°C | Давл.: %.0f Pa\n",
		env.WeatherType, env.Temperature, env.Pressure)

	params := []struct{ k, v string }{
		{"Weibull k", fmt.Sprintf("%.2f", env.WeibullK)},
		{"Weibull c", fmt.Sprintf("%.2f м/с", env.WeibullC)},
		{"Roughness", fmt.Sprintf("%.2f", loc.Roughness)},
		{"HubHeight", fmt.Sprintf("%.1f м", loc.HubHeight)},
		{"BaseHeight", fmt.Sprintf("%.1f м", loc.BaseHeight)},
	}
	fmt.Print("Параметры ветра:")
	for _, p := range params {
		fmt.Printf(" • %s=%s", p.k, p.v)
	}
	fmt.Println()
	fmt.Println("Готово: можно рассчитывать выработку энергии.")
}
