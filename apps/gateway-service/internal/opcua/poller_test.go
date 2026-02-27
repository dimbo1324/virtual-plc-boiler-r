package opcua

import (
	"context"
	"errors"
	"gateway-service/internal/domain"
	"testing"

	"github.com/gopcua/opcua/ua"
	"github.com/stretchr/testify/assert"
)

type mockOpcuaClient struct {
	connectErr error
	readResp   *ua.ReadResponse
	readErr    error
	closeErr   error
}

func (m *mockOpcuaClient) Connect(ctx context.Context) error { return m.connectErr }
func (m *mockOpcuaClient) Read(ctx context.Context, req *ua.ReadRequest) (*ua.ReadResponse, error) {
	return m.readResp, m.readErr
}
func (m *mockOpcuaClient) Close(ctx context.Context) error { return m.closeErr }

func TestOpcClientConnectSuccess(t *testing.T) {
	client := NewOpcClient("test")
	client.client = &mockOpcuaClient{}
	err := client.Connect(context.Background())
	assert.NoError(t, err)
}

func TestOpcClientConnectRetryFail(t *testing.T) {
	client := NewOpcClient("test")
	err := client.Connect(context.Background())
	assert.Error(t, err)
}

func TestOpcClientReadSuccess(t *testing.T) {
	client := NewOpcClient("test")
	client.client = &mockOpcuaClient{
		readResp: &ua.ReadResponse{
			Results: []*ua.DataValue{
				{Value: ua.MustVariant(float64(450.5))},
				{Value: ua.MustVariant(float64(60.2))},
			},
		},
	}
	tags, err := client.Read(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, 450.5, tags.Temperature)
	assert.Equal(t, 60.2, tags.Pressure)
}

func TestOpcClientReadError(t *testing.T) {
	client := NewOpcClient("test")
	client.client = &mockOpcuaClient{readErr: errors.New("read fail")}
	tags, err := client.Read(context.Background())
	assert.Error(t, err)
	assert.Equal(t, domain.Tags{}, tags)
}

func TestOpcClientReadNilClient(t *testing.T) {
	client := NewOpcClient("test")
	tags, err := client.Read(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, domain.Tags{}, tags)
}

func TestOpcClientClose(t *testing.T) {
	client := NewOpcClient("test")
	client.client = &mockOpcuaClient{}
	assert.NoError(t, client.Close())
}

func TestOpcClientCloseNil(t *testing.T) {
	client := NewOpcClient("test")
	assert.NoError(t, client.Close())
}

func TestMockPoller(t *testing.T) {
	poller := NewMockPoller()
	assert.NoError(t, poller.Connect(context.Background()))
	tags, err := poller.Read(context.Background())
	assert.NoError(t, err)
	assert.Greater(t, tags.Temperature, 400.0)
	assert.Less(t, tags.Temperature, 520.0)
	assert.Greater(t, tags.Pressure, 55.0)
	assert.Less(t, tags.Pressure, 75.0)
	assert.NoError(t, poller.Close())
}
