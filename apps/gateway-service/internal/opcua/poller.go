package opcua

import (
	"context"
	"fmt"
	"gateway-service/internal/domain"

	"github.com/gopcua/opcua"
	"github.com/gopcua/opcua/ua"
	"go.uber.org/zap"
)

type IPoller interface {
	Connect(ctx context.Context) error
	Read(ctx context.Context) (domain.Tags, error)
	Close() error
}
type OpcClient struct {
	endpoint string
	client   *opcua.Client
	logger   *zap.SugaredLogger
	nsIndex  uint16
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
	c.client, err = opcua.NewClient(c.endpoint, opcua.SecurityMode(ua.MessageSecurityModeNone))
	if err != nil {
		return err
	}
	if err = c.client.Connect(ctx); err != nil {
		return err
	}
	namespaceURI := "urn:virtual-plc:boiler"
	namespaces := c.client.Namespaces()
	for i, uri := range namespaces {
		if uri == namespaceURI {
			c.nsIndex = uint16(i)
			c.logger.Infof("Found namespace '%s' at index %d", namespaceURI, i)
			break
		}
	}
	if c.nsIndex == 0 {
		c.logger.Warnf("Namespace '%s' not found, defaulting to index 2", namespaceURI)
		c.nsIndex = 2
	}
	return nil
}
func (c *OpcClient) Read(ctx context.Context) (domain.Tags, error) {
	if c.client == nil {
		return domain.Tags{}, fmt.Errorf("client not connected")
	}
	req := &ua.ReadRequest{
		NodesToRead: []*ua.ReadValueID{
			{NodeID: ua.NewNumericNodeID(c.nsIndex, 1002), AttributeID: ua.AttributeIDValue},
			{NodeID: ua.NewNumericNodeID(c.nsIndex, 1001), AttributeID: ua.AttributeIDValue},
			{NodeID: ua.NewNumericNodeID(c.nsIndex, 1003), AttributeID: ua.AttributeIDValue},
			{NodeID: ua.NewNumericNodeID(c.nsIndex, 1005), AttributeID: ua.AttributeIDValue},
			{NodeID: ua.NewNumericNodeID(c.nsIndex, 1006), AttributeID: ua.AttributeIDValue},
		},
	}
	resp, err := c.client.Read(ctx, req)
	if err != nil {
		return domain.Tags{}, err
	}
	if resp.Results[0].Status != ua.StatusOK {
		return domain.Tags{}, fmt.Errorf("node read failed with status: %v", resp.Results[0].Status)
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

func NewMockPoller() *MockPoller                        { return &MockPoller{} }
func (m *MockPoller) Connect(ctx context.Context) error { return nil }
func (m *MockPoller) Close() error                      { return nil }
func (m *MockPoller) Read(ctx context.Context) (domain.Tags, error) {
	return domain.Tags{Temperature: 100, Pressure: 10}, nil
}
