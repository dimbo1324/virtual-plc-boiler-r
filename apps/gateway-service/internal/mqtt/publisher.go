package mqtt

import (
	"gateway-service/internal/domain"
	"log"
)

type IPublisher interface {
	Connect() error
	Publish(topic string, payload domain.Payload) error
	Close()
}

type MqttClient struct {
	BrokerURL string
	ClientID  string
}

func NewMqttClient(brokerURL, clientID string) *MqttClient {
	return &MqttClient{BrokerURL: brokerURL, ClientID: clientID}
}

func (m *MqttClient) Connect() error                                     { return nil }
func (m *MqttClient) Publish(topic string, payload domain.Payload) error { return nil }
func (m *MqttClient) Close()                                             {}

type MockPublisher struct{}

func (m *MockPublisher) Connect() error { return nil }
func (m *MockPublisher) Publish(topic string, payload domain.Payload) error {
	log.Printf("[MOCK MQTT] Topic: %s | Payload: %+v", topic, payload)
	return nil
}

func (m *MockPublisher) Close() {
	log.Println("MQTT: Closing connection...")
}
