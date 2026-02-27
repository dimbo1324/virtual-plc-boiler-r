package opcua

import (
	"context"
	"gateway-service/internal/domain"
	"math/rand"
)

type IPoller interface {
	Connect(ctx context.Context) error
	Read(ctx context.Context) (domain.Tags, error)
	Close() error
}

type OpcClient struct {
	Endpoint string
}

func NewOpcClient(endpoint string) *OpcClient {
	return &OpcClient{Endpoint: endpoint}
}

func (c *OpcClient) Connect(ctx context.Context) error { return nil }
func (c *OpcClient) Close() error                      { return nil }
func (c *OpcClient) Read(ctx context.Context) (domain.Tags, error) {
	return domain.Tags{}, nil
}

type MockPoller struct{}

func (m *MockPoller) Connect(ctx context.Context) error { return nil }
func (m *MockPoller) Close() error                      { return nil }
func (m *MockPoller) Read(ctx context.Context) (domain.Tags, error) {
	return domain.Tags{
		Temperature: 400.0 + rand.Float64()*100,
		Pressure:    55.0 + rand.Float64()*10,
	}, nil
}

// type OpcClient struct {
// 	Endpoint string
// }

// func NewOpcClient(endpoint string) *OpcClient {
// 	return &OpcClient{Endpoint: endpoint}
// }
