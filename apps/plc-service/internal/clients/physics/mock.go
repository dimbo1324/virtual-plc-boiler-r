package physics

import (
	"context"
	pb "plc-service/pkg/grpc"

	"github.com/stretchr/testify/mock"
)

type MockPhysicsClient struct {
	mock.Mock
}

func (m *MockPhysicsClient) GetStatus(ctx context.Context) (*pb.BoilerStatus, error) {
	args := m.Called(ctx)
	if args.Get(0) != nil {
		return args.Get(0).(*pb.BoilerStatus), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockPhysicsClient) SetControls(ctx context.Context, fuel, water, steam float64) (*pb.BoilerStatus, error) {
	args := m.Called(ctx, fuel, water, steam)
	if args.Get(0) != nil {
		return args.Get(0).(*pb.BoilerStatus), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockPhysicsClient) Close() error {
	args := m.Called()
	return args.Error(0)
}
