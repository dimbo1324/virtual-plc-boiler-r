package worker

import (
	"context"
	"gateway-service/internal/domain"
	"gateway-service/internal/mqtt"
	"sync"
	"time"

	"go.uber.org/zap"
)

type Pool struct {
	jobsChan   chan domain.Payload
	retryChan  chan domain.Payload
	publisher  mqtt.IPublisher
	wg         sync.WaitGroup
	logger     *zap.SugaredLogger
	maxRetries int
}

func NewPool(publisher mqtt.IPublisher, bufferSize int, logger *zap.Logger) *Pool {
	return &Pool{
		jobsChan:   make(chan domain.Payload, bufferSize),
		retryChan:  make(chan domain.Payload, 200),
		publisher:  publisher,
		logger:     logger.Sugar(),
		maxRetries: 3,
	}
}

func (p *Pool) Start(ctx context.Context, workersCount int) {
	for i := 0; i < workersCount; i++ {
		p.wg.Add(1)
		go p.worker(ctx, i)
	}
	go p.retryWorker(ctx)
	p.logger.Infof("Worker pool started with %d workers", workersCount)
}

func (p *Pool) worker(ctx context.Context, id int) {
	defer p.wg.Done()
	for {
		select {
		case <-ctx.Done():
			return
		case job, ok := <-p.jobsChan:
			if !ok {
				return
			}
			p.publishWithRetry(job, 0)
		}
	}
}

func (p *Pool) retryWorker(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case job := <-p.retryChan:
			p.publishWithRetry(job, 1)
		}
	}
}

func (p *Pool) publishWithRetry(job domain.Payload, attempt int) {
	err := p.publisher.Publish("v1/devices/boiler/telemetry", job)
	if err == nil {
		return
	}

	if attempt < p.maxRetries {
		time.Sleep(time.Duration(1<<attempt) * 300 * time.Millisecond)
		p.retryChan <- job
		p.logger.Warnw("Publish failed, moved to retry queue", "attempt", attempt, "err", err)
	} else {
		p.logger.Errorw("Publish failed after max retries, dropping", "err", err)
	}
}

func (p *Pool) Push(payload domain.Payload) {
	select {
	case p.jobsChan <- payload:
	default:
		p.retryChan <- payload
		p.logger.Warn("Main queue full → moved to retry queue")
	}
}

func (p *Pool) Stop() {
	close(p.jobsChan)
	close(p.retryChan)
	p.wg.Wait()
	p.logger.Info("Worker pool stopped")
}
