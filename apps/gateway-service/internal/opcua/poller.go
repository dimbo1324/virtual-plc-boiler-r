package opcua

import (
	"context"
	"gateway-service/internal/domain"
)

type IPoller interface {
	Connect(ctx context.Context) error
	Read() (domain.Tags, error)
	Close() error
}

type OpcClient struct {
	Endpoint string
}

func NewOpcClient(endpoint string) *OpcClient {
	return &OpcClient{Endpoint: endpoint}
}
