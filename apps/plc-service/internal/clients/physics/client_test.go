package physics

import (
	"context"
	"testing"

	pb "plc-service/pkg/grpc"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestMockPhysicsClient_Interaction(t *testing.T) {
	mockClient := new(MockPhysicsClient)

	mockClient.On("GetStatus", mock.Anything).Return(&pb.BoilerStatus{
		Timestamp:     10.0,
		SteamPressure: 50.0,
		FurnaceTemp:   400.0,
	}, nil)

	mockClient.On("SetControls", mock.Anything, 100.0, 50.0, 0.0).Return(&pb.BoilerStatus{}, nil)

	ctx := context.Background()

	status, err := mockClient.GetStatus(ctx)
	assert.NoError(t, err)
	assert.Equal(t, 50.0, status.SteamPressure)

	_, err = mockClient.SetControls(ctx, 100.0, 50.0, 0.0)
	assert.NoError(t, err)

	mockClient.AssertExpectations(t)
}
