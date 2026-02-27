package mqtt

import (
	"gateway-service/internal/domain"
	"testing"
	"time"

	paho "github.com/eclipse/paho.mqtt.golang"
	"github.com/stretchr/testify/assert"
)

type mockToken struct {
	err error
}

func (t *mockToken) Wait() bool                       { return true }
func (t *mockToken) WaitTimeout(d time.Duration) bool { return true }
func (t *mockToken) Error() error                     { return t.err }
func (t *mockToken) Done() <-chan struct{}            { ch := make(chan struct{}); close(ch); return ch }

type mockPahoClient struct {
	connectErr error
	publishErr error
	connected  bool
}

func (m *mockPahoClient) IsConnected() bool      { return m.connected }
func (m *mockPahoClient) IsConnectionOpen() bool { return m.connected }

func (m *mockPahoClient) Connect() paho.Token {
	m.connected = true
	return &mockToken{err: m.connectErr}
}

func (m *mockPahoClient) Publish(topic string, qos byte, retained bool, payload interface{}) paho.Token {
	return &mockToken{err: m.publishErr}
}

func (m *mockPahoClient) Disconnect(quiesce uint) {}

func (m *mockPahoClient) AddRoute(topic string, callback paho.MessageHandler) {}
func (m *mockPahoClient) Subscribe(topic string, qos byte, callback paho.MessageHandler) paho.Token {
	return &mockToken{}
}
func (m *mockPahoClient) SubscribeMultiple(filters map[string]byte, callback paho.MessageHandler) paho.Token {
	return &mockToken{}
}
func (m *mockPahoClient) Unsubscribe(topics ...string) paho.Token {
	return &mockToken{}
}
func (m *mockPahoClient) OptionsReader() paho.ClientOptionsReader {
	return paho.ClientOptionsReader{}
}

func TestMqttClientConnect(t *testing.T) {
	client := NewMqttClient("tcp://test:1883", "test_id")
	client.client = &mockPahoClient{}
	err := client.Connect()
	assert.NoError(t, err)
}

func TestMqttClientPublishSuccess(t *testing.T) {
	client := NewMqttClient("tcp://test:1883", "test_id")
	client.client = &mockPahoClient{}

	payload := domain.Payload{
		Timestamp: "2026-02-27T10:00:00Z",
		AssetID:   "boiler_01",
		Tags: domain.Tags{
			Temperature: 450.5,
			Pressure:    60.2,
		},
	}

	err := client.Publish("v1/devices/boiler/telemetry", payload)
	assert.NoError(t, err)
}

func TestMqttClientPublishError(t *testing.T) {
	client := NewMqttClient("tcp://test:1883", "test_id")
	client.client = &mockPahoClient{publishErr: assert.AnError}

	payload := domain.Payload{Timestamp: "time", AssetID: "id", Tags: domain.Tags{Temperature: 1, Pressure: 2}}
	err := client.Publish("topic", payload)
	assert.Error(t, err)
}

func TestMockPublisher(t *testing.T) {
	pub := NewMockPublisher()
	assert.NoError(t, pub.Connect())

	payload := domain.Payload{Timestamp: "time", AssetID: "id", Tags: domain.Tags{Temperature: 1, Pressure: 2}}
	assert.NoError(t, pub.Publish("topic", payload))

	pub.Close()
}
