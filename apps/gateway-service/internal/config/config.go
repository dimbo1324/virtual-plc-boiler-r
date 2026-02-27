package config

import "github.com/spf13/viper"

type Config struct {
	OPCUAEndpoint  string
	MQTTBroker     string
	MQTTClientID   string
	Topic          string
	UseMocks       bool
	WorkerCount    int
	BufferSize     int
	PollIntervalMs int
}

func Load() *Config {
	viper.SetDefault("opcua.endpoint", "opc.tcp://localhost:4840")
	viper.SetDefault("mqtt.broker", "tcp://localhost:1883")
	viper.SetDefault("mqtt.client_id", "gw_01")
	viper.SetDefault("mqtt.topic", "v1/devices/boiler/telemetry")
	viper.SetDefault("use_mocks", true)
	viper.SetDefault("worker.count", 5)
	viper.SetDefault("buffer.size", 500)
	viper.SetDefault("poll.interval_ms", 500)
	viper.AutomaticEnv()
	return &Config{
		OPCUAEndpoint:  viper.GetString("opcua.endpoint"),
		MQTTBroker:     viper.GetString("mqtt.broker"),
		MQTTClientID:   viper.GetString("mqtt.client_id"),
		Topic:          viper.GetString("mqtt.topic"),
		UseMocks:       viper.GetBool("use_mocks"),
		WorkerCount:    viper.GetInt("worker.count"),
		BufferSize:     viper.GetInt("buffer.size"),
		PollIntervalMs: viper.GetInt("poll.interval_ms"),
	}
}
