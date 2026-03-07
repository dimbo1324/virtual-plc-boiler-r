# 🔥 Virtual PLC Boiler

A microservice-based digital twin of an industrial steam boiler, originally developed as a proof-of-concept for a virtual power plant (VPP) project at the **Jazan Industrial Complex, Saudi Arabia**.

The system simulates a real PLC-controlled boiler loop end-to-end: from physical thermodynamics all the way to a live Grafana dashboard — using nothing but Docker and open-source tech.

---

## How It Works

At its core, the project is a chain of five cooperating services that mimic a real industrial telemetry pipeline.

```
┌─────────────────────────────────────────────────────────────────┐
│                                                                 │
│  ┌──────────────────┐   gRPC    ┌──────────────────┐           │
│  │  physics-service │ ◄───────► │   plc-service    │           │
│  │   (Python)       │           │     (Go)         │           │
│  │                  │           │                  │           │
│  │  Boiler physics  │           │  PID controllers │           │
│  │  Steam tables    │           │  OPC-UA server   │           │
│  └──────────────────┘           └────────┬─────────┘           │
│                                          │ OPC-UA               │
│                                 ┌────────▼─────────┐           │
│                                 │ gateway-service  │           │
│                                 │     (Go)         │           │
│                                 │  Worker pool     │           │
│                                 └────────┬─────────┘           │
│                                          │ MQTT                 │
│                                 ┌────────▼─────────┐           │
│                                 │    mosquitto     │           │
│                                 │  (MQTT broker)   │           │
│                                 └────────┬─────────┘           │
│                                          │ MQTT                 │
│                                 ┌────────▼─────────┐           │
│                          ┌──────┤data-consumer-svc │           │
│                          │      │    (Python)      │           │
│                          │      └──────────────────┘           │
│                    write │                                      │
│                  ┌───────▼──────┐    query   ┌──────────────┐  │
│                  │   InfluxDB   │ ◄───────── │   Grafana    │  │
│                  │  (TSDB)      │            │  Dashboard   │  │
│                  └──────────────┘            └──────────────┘  │
└─────────────────────────────────────────────────────────────────┘
```

### Step by step

**1. physics-service** is the "real world". It runs a continuous simulation of boiler thermodynamics: furnace temperature, steam pressure (derived from a steam table via linear interpolation), drum water level, and steam flow. It exposes a gRPC API with two methods — `GetStatus` to read sensor values and `SetControls` to adjust the three control valves (fuel, feedwater, steam).

**2. plc-service** acts as the virtual PLC. Every second it queries the physics service for the current boiler state, then runs two independent PID controllers — one to regulate steam pressure via the fuel valve, another to maintain the drum water level via the feedwater valve. The computed commands are sent back to the physics service via gRPC. Simultaneously, the PLC exposes an OPC-UA server (port 4840) broadcasting the live telemetry so any SCADA-compatible client can connect.

**3. gateway-service** is the IIoT edge gateway. It polls the OPC-UA server every 500 ms using a configurable worker pool, serializes the readings into a JSON payload tagged with `asset_id: boiler_01`, and publishes them to the Mosquitto MQTT broker on topic `v1/devices/boiler/telemetry`.

**4. data-consumer-service** subscribes to that MQTT topic and writes each incoming message as a `boiler_telemetry` measurement into InfluxDB (bucket: `telemetry`), including temperature, pressure, fuel position, drum level, and steam flow.

**5. Grafana** connects to InfluxDB and renders the live dashboard with pre-provisioned charts for all key process variables.

---

## Services

| Service                 | Language    | Role                             | Port            |
| ----------------------- | ----------- | -------------------------------- | --------------- |
| `physics-service`       | Python 3.12 | Boiler physics simulation        | `50051` (gRPC)  |
| `plc-service`           | Go          | PID control loop + OPC-UA server | `4840` (OPC-UA) |
| `gateway-service`       | Go          | OPC-UA → MQTT bridge             | —               |
| `mosquitto`             | —           | MQTT broker                      | `1883`, `9001`  |
| `data-consumer-service` | Python 3.12 | MQTT → InfluxDB writer           | —               |
| `influxdb`              | —           | Time-series database             | `8086`          |
| `grafana`               | —           | Visualization dashboard          | `3001`          |

---

## Tech Stack

- **Simulation** — Python, dataclasses, interpolated steam tables (0–1500 °C / 0–500 bar)
- **Control** — Go, dual PID controllers (pressure loop + drum level loop)
- **Industrial protocols** — gRPC (PLC ↔ Physics), OPC-UA (PLC → Gateway), MQTT (Gateway → Broker)
- **Observability** — InfluxDB 2, Grafana, Zap structured logging
- **Infrastructure** — Docker Compose, multi-stage builds, healthchecks

---

## Quick Start

```bash
git clone https://github.com/your-username/virtual-plc-boiler.git
cd virtual-plc-boiler/apps

docker compose up --build
```

Open Grafana at [http://localhost:3001](http://localhost:3001) — credentials `admin / admin`.

> The startup order is enforced via Docker healthchecks:
> `physics-service` must be healthy before `plc-service` starts,
> which in turn must be healthy before `gateway-service` connects.

---

## Running Tests

**Go services:**
```bash
# gateway-service
cd apps/gateway-service && go test ./... -v -cover

# plc-service
cd apps/plc-service && go test ./... -v -cover
```

**Python (physics-service):**
```bash
cd apps/physics-service
uv run pytest tests/ -v
```

---

## Physics Model

The simulation uses a simplified but physically grounded model:

- **Temperature** follows a first-order lag toward a target set by the fuel valve opening, with configurable heating (`0.05 /s`) and cooling (`0.02 /s`) rates.
- **Steam pressure** is derived from temperature via a piecewise-linear steam table spanning 0–1500 °C, then reduced by the steam valve opening.
- **Drum level** rises with feedwater inflow (`2.5 mm/s at 100%`) and falls due to evaporation proportional to furnace temperature.

All parameters are tunable via environment variables (e.g. `BOILER_HEATING_RATE`, `BOILER_MAX_FURNACE_TEMP`).

---

## Background

This project was prototyped as part of a virtual power plant (VPP) platform study for an energy boiler at the **Jazan Economic City industrial complex** in Saudi Arabia. The goal was to validate the full telemetry pipeline — from physics to dashboard — in a fully containerized environment before any integration with physical hardware.

---

## License

MIT
