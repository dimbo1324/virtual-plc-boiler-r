package mqtt

import (
	"encoding/json"
	"gateway-service/internal/domain"
	"log"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"go.uber.org/zap"
)

type IPublisher interface {
	Connect() error
	Publish(topic string, payload domain.Payload) error
	Close()
}

type MqttClient struct {
	client mqtt.Client
	logger *zap.SugaredLogger
}

func NewMqttClient(broker, clientID string) *MqttClient {
	opts := mqtt.NewClientOptions().
		AddBroker(broker).
		SetClientID(clientID).
		SetAutoReconnect(true).
		SetMaxReconnectInterval(5 * time.Second)

	logger, _ := zap.NewProduction()

	return &MqttClient{
		client: mqtt.NewClient(opts),
		logger: logger.Sugar(),
	}
}

func (m *MqttClient) Connect() error {
	token := m.client.Connect()
	token.Wait()
	if err := token.Error(); err != nil {
		m.logger.Errorw("MQTT connection failed", "err", err)
		return err
	}
	m.logger.Info("MQTT connected")
	return nil
}

func (m *MqttClient) Publish(topic string, payload domain.Payload) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	token := m.client.Publish(topic, 1, false, data)
	token.Wait()
	return token.Error()
}

func (m *MqttClient) Close() {
	m.client.Disconnect(250)
}

type MockPublisher struct{}

func NewMockPublisher() *MockPublisher {
	return &MockPublisher{}
}

func (m *MockPublisher) Connect() error {
	log.Println("[MOCK MQTT] Connected")
	return nil
}

func (m *MockPublisher) Publish(topic string, payload domain.Payload) error {
	log.Printf("[MOCK MQTT] Topic: %s | Payload: %+v", topic, payload)
	return nil
}

func (m *MockPublisher) Close() {
	log.Println("[MOCK MQTT] Closing connection...")
}
