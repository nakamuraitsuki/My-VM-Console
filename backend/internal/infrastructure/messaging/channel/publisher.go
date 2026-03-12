package channel

import (
	"context"
	"encoding/json"
	"fmt"

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

func (p *publisher) Publish(ctx context.Context, jobType usecase.JobType, payload any) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	p.hub.mu.RLock()
	defer p.hub.mu.RUnlock()

	channels, ok := p.hub.channels[string(jobType)]
	if !ok {
		return nil // 購読者がいなければ何もしない
	}

	// 全ての購読者にブロードキャスト
	for _, ch := range channels {
		go func(c chan []byte) {
			select {
			case c <- data:
			case <-ctx.Done():
			}
		}(ch)
	}

	return nil
}
