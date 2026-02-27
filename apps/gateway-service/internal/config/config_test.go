package config

import (
	"os"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestLoad(t *testing.T) {
	cfg := Load()
	assert.Equal(t, "opc.tcp://localhost:4840", cfg.OPCUAEndpoint)
	assert.Equal(t, "tcp://localhost:1883", cfg.MQTTBroker)
	assert.Equal(t, "gw_01", cfg.MQTTClientID)
	assert.Equal(t, "v1/devices/boiler/telemetry", cfg.Topic)
	assert.True(t, cfg.UseMocks)
	assert.Equal(t, 5, cfg.WorkerCount)
	assert.Equal(t, 500, cfg.BufferSize)
	assert.Equal(t, 500, cfg.PollIntervalMs)

	os.Setenv("OPCUA_ENDPOINT", "opc.tcp://test:1234")
	os.Setenv("MQTT_BROKER", "tcp://broker:9999")
	os.Setenv("MQTT_CLIENT_ID", "test_id")
	os.Setenv("MQTT_TOPIC", "test/topic")
	os.Setenv("USE_MOCKS", "false")
	os.Setenv("WORKER_COUNT", "10")
	os.Setenv("BUFFER_SIZE", "1000")
	os.Setenv("POLL_INTERVAL_MS", "100")

	viper.Reset()
	cfg = Load()
	assert.Equal(t, "opc.tcp://test:1234", cfg.OPCUAEndpoint)
	assert.Equal(t, "tcp://broker:9999", cfg.MQTTBroker)
	assert.Equal(t, "test_id", cfg.MQTTClientID)
	assert.Equal(t, "test/topic", cfg.Topic)
	assert.False(t, cfg.UseMocks)
	assert.Equal(t, 10, cfg.WorkerCount)
	assert.Equal(t, 1000, cfg.BufferSize)
	assert.Equal(t, 100, cfg.PollIntervalMs)

	os.Unsetenv("OPCUA_ENDPOINT")
	os.Unsetenv("MQTT_BROKER")
	os.Unsetenv("MQTT_CLIENT_ID")
	os.Unsetenv("MQTT_TOPIC")
	os.Unsetenv("USE_MOCKS")
	os.Unsetenv("WORKER_COUNT")
	os.Unsetenv("BUFFER_SIZE")
	os.Unsetenv("POLL_INTERVAL_MS")
}
