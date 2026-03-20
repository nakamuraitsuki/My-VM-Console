package channel

import (
	"context"
	"log"

	"example.com/m/internal/usecase"
)

type publisher struct {
	hub *Hub
}

func NewPublisher(hub *Hub) usecase.JobPublisher {
	return &publisher{
		hub: hub,
	}
}

func (p *publisher) Publish(ctx context.Context, jobType usecase.JobType, payload []byte) error {
	p.hub.mu.RLock()
	defer p.hub.mu.RUnlock()

	log.Printf("[Publisher] publishing job: type=%s", jobType)
	channels, ok := p.hub.channels[string(jobType)]
	if !ok {
		log.Printf("[Publisher] no subscribers for job type: %s", jobType)
		return nil // 購読者がいなければ何もしない
	}

	// 全ての購読者にブロードキャスト
	for _, ch := range channels {
		log.Printf("[Publisher] sending job to subscriber: type=%s", jobType)
		select {
		case ch <- payload:
		case <-ctx.Done():
			log.Printf("[Publisher] context cancelled while sending job: type=%s", jobType)
		}
	}

	return nil
}
