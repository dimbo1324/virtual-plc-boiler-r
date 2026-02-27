package worker

import (
	"context"
	"encoding/json"
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
	log.Printf("Worker %d: ready", id)

	for {
		select {
		case <-ctx.Done():
			log.Printf("Worker %d stopping...", id)
			return
		case job, ok := <-p.jobsChan:
			if !ok {
				return
			}

			_, err := json.Marshal(job)
			if err != nil {
				log.Printf("Worker %d: error marshalling: %v", id, err)
				continue
			}

			err = p.publisher.Publish("v1/devices/boiler/telemetry", job)
			if err != nil {
				log.Printf("Worker %d: publish error: %v", id, err)
			}
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
	close(p.jobsChan)
	p.wg.Wait()
	log.Println("Worker pool stopped cleanly")
}
