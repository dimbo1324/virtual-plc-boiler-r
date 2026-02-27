package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoad_Defaults(t *testing.T) {
	cfg := Load()

	assert.Equal(t, "opc.tcp://localhost:4840", cfg.OPCUAEndpoint)
	assert.Equal(t, "tcp://localhost:1883", cfg.MQTTBroker)
	assert.Equal(t, "gw_01", cfg.MQTTClientID)
	assert.Equal(t, "v1/devices/boiler/telemetry", cfg.Topic)
	assert.True(t, cfg.UseMocks)
	assert.Equal(t, 5, cfg.WorkerCount)
	assert.Equal(t, 500, cfg.BufferSize)
	assert.Equal(t, 500, cfg.PollIntervalMs)
}

func TestLoad_FromEnvironment(t *testing.T) {
	t.Setenv("OPCUA_ENDPOINT", "opc.tcp://test:1234")
	t.Setenv("MQTT_BROKER", "tcp://broker:9999")
	t.Setenv("MQTT_CLIENT_ID", "test_id")
	t.Setenv("MQTT_TOPIC", "test/topic")
	t.Setenv("USE_MOCKS", "false")
	t.Setenv("WORKER_COUNT", "10")
	t.Setenv("BUFFER_SIZE", "1000")
	t.Setenv("POLL_INTERVAL_MS", "100")

	cfg := Load()

	assert.Equal(t, "opc.tcp://test:1234", cfg.OPCUAEndpoint)
	assert.Equal(t, "tcp://broker:9999", cfg.MQTTBroker)
	assert.Equal(t, "test_id", cfg.MQTTClientID)
	assert.Equal(t, "test/topic", cfg.Topic)
	assert.False(t, cfg.UseMocks)
	assert.Equal(t, 10, cfg.WorkerCount)
	assert.Equal(t, 1000, cfg.BufferSize)
	assert.Equal(t, 100, cfg.PollIntervalMs)
}
