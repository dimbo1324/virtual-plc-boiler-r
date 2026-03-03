package opcua

import (
	"context"
	"gateway-service/internal/domain"
	"math/rand"
	"time"

	"github.com/gopcua/opcua"
	"github.com/gopcua/opcua/ua"
	"go.uber.org/zap"
)

type IPoller interface {
	Connect(ctx context.Context) error
	Read(ctx context.Context) (domain.Tags, error)
	Close() error
}

type opcuaClient interface {
	Connect(ctx context.Context) error
	Read(ctx context.Context, req *ua.ReadRequest) (*ua.ReadResponse, error)
	Close(ctx context.Context) error
}

type OpcClient struct {
	endpoint string
	client   opcuaClient
	logger   *zap.SugaredLogger
}

func NewOpcClient(endpoint string) *OpcClient {
	logger, _ := zap.NewProduction()
	return &OpcClient{
		endpoint: endpoint,
		logger:   logger.Sugar(),
	}
}

func (c *OpcClient) Connect(ctx context.Context) error {
	var err error

	for i := 0; i < 5; i++ {
		c.client, err = opcua.NewClient(c.endpoint)
		if err != nil {
			c.logger.Warnw("Failed to create OPC UA client", "attempt", i+1, "err", err)
			time.Sleep(time.Duration(1<<uint(i)) * time.Second)
			continue
		}

		if err = c.client.Connect(ctx); err == nil {
			c.logger.Info("OPC UA connected successfully")
			return nil
		}

		c.logger.Warnw("OPC UA connection failed, retrying...", "attempt", i+1, "err", err)
		time.Sleep(time.Duration(1<<uint(i)) * time.Second)
	}

	c.logger.Error("All OPC UA connection attempts failed")
	return err
}

func (c *OpcClient) Read(ctx context.Context) (domain.Tags, error) {
	if c.client == nil {
		return domain.Tags{}, nil
	}

	req := &ua.ReadRequest{
		NodesToRead: []*ua.ReadValueID{
			{NodeID: ua.NewNumericNodeID(2, 1002)},
			{NodeID: ua.NewNumericNodeID(2, 1001)},
			{NodeID: ua.NewNumericNodeID(2, 1003)},
			{NodeID: ua.NewNumericNodeID(2, 1005)},
			{NodeID: ua.NewNumericNodeID(2, 1006)},
		},
	}

	resp, err := c.client.Read(ctx, req)
	if err != nil {
		c.logger.Warnw("OPC UA read failed", "err", err)
		return domain.Tags{}, err
	}

	return domain.Tags{
		Temperature: resp.Results[0].Value.Float(),
		Pressure:    resp.Results[1].Value.Float(),
		Fuel:        resp.Results[2].Value.Float(),
		DrumLevel:   resp.Results[3].Value.Float(),
		SteamFlow:   resp.Results[4].Value.Float(),
	}, nil
}

func (c *OpcClient) Close() error {
	if c.client != nil {
		return c.client.Close(context.Background())
	}
	return nil
}

type MockPoller struct{}

func NewMockPoller() *MockPoller {
	return &MockPoller{}
}

func (m *MockPoller) Connect(ctx context.Context) error { return nil }
func (m *MockPoller) Close() error                      { return nil }

func (m *MockPoller) Read(ctx context.Context) (domain.Tags, error) {
	return domain.Tags{
		Temperature: 400.0 + rand.Float64()*120,
		Pressure:    55.0 + rand.Float64()*20,
	}, nil
}
