package mqtt

import (
	"gateway-service/internal/domain"
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
