package worker

import (
	"context"
	"log"
	"sync"

	"gateway-service/internal/domain"
	"gateway-service/internal/mqtt"
)

type Pool struct {
	jobsChan  chan domain.Payload
	publisher mqtt.IPublisher
	wg        sync.WaitGroup
}

func NewPool(publisher mqtt.IPublisher, bufferSize int) *Pool {
	return &Pool{
		jobsChan:  make(chan domain.Payload, bufferSize),
		publisher: publisher,
	}
}

func (p *Pool) Start(ctx context.Context, workersCount int) {
	for i := 0; i < workersCount; i++ {
		p.wg.Add(1)
		go p.worker(ctx, i)
	}
	log.Printf("Worker pool started with %d workers", workersCount)
}

func (p *Pool) worker(ctx context.Context, id int) {
	defer p.wg.Done()
	for {
		select {
		case <-ctx.Done():
			log.Printf("Worker %d stopping...", id)
			return
		case job := <-p.jobsChan:
			_ = job
		}
	}
}

func (p *Pool) Push(payload domain.Payload) {
	select {
	case p.jobsChan <- payload:
	default:
		log.Println("Warning: Job queue is full! Dropping payload.")
	}
}

func (p *Pool) Stop() {
	p.wg.Wait()
}
