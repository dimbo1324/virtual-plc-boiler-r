package worker

import (
	"context"
	"errors"
	"gateway-service/internal/domain"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

type mockPublisher struct {
	publishErr error
	published  []domain.Payload
}

func (m *mockPublisher) Connect() error { return nil }
func (m *mockPublisher) Publish(topic string, payload domain.Payload) error {
	if m.publishErr == nil {
		m.published = append(m.published, payload)
	}
	return m.publishErr
}
func (m *mockPublisher) Close() {}

func TestPoolStartStop(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	pub := &mockPublisher{}
	pool := NewPool(pub, 10, logger)
	ctx, cancel := context.WithCancel(context.Background())
	pool.Start(ctx, 2)
	time.Sleep(100 * time.Millisecond)
	cancel()
	pool.Stop()
}

func TestPoolPushAndPublishSuccess(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	pub := &mockPublisher{}
	pool := NewPool(pub, 10, logger)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	pool.Start(ctx, 2)

	payload := domain.Payload{Timestamp: "time", AssetID: "id", Tags: domain.Tags{Temperature: 1.0, Pressure: 2.0}}
	pool.Push(payload)
	time.Sleep(200 * time.Millisecond)
	assert.Len(t, pub.published, 1)
	assert.Equal(t, payload, pub.published[0])
}

func TestPoolRetryOnFail(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	pub := &mockPublisher{publishErr: errors.New("fail")}
	pool := NewPool(pub, 10, logger)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	pool.Start(ctx, 2)

	payload := domain.Payload{Timestamp: "time", AssetID: "id", Tags: domain.Tags{Temperature: 1.0, Pressure: 2.0}}
	pool.Push(payload)
	time.Sleep(2 * time.Second)
	assert.Len(t, pub.published, 0)
}

func TestPoolBackpressureToRetry(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	pub := &mockPublisher{}
	pool := NewPool(pub, 1, logger)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	pool.Start(ctx, 1)

	pool.Push(domain.Payload{})
	pool.Push(domain.Payload{})

	time.Sleep(200 * time.Millisecond)
	assert.Len(t, pub.published, 2)
}

func TestPoolPublishWithRetrySuccessAfterFail(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	pub := &mockPublisher{publishErr: errors.New("fail")}
	pool := NewPool(pub, 10, logger)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	pool.Start(ctx, 2)

	payload := domain.Payload{Timestamp: "time", AssetID: "id", Tags: domain.Tags{Temperature: 1.0, Pressure: 2.0}}
	pool.Push(payload)
	time.Sleep(500 * time.Millisecond)

	pub.publishErr = nil
	time.Sleep(500 * time.Millisecond)
	assert.Len(t, pub.published, 1)
}
