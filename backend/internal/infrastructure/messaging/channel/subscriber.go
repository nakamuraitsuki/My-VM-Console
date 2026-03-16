package channel

import (
	"context"
	"log"

	"example.com/m/internal/usecase"
)

type subscriber struct {
	hub *Hub
}

func NewSubscriber(hub *Hub) usecase.JobSubscriber {
	return &subscriber{
		hub: hub,
	}
}

func (s *subscriber) Subscribe(ctx context.Context, jobType usecase.JobType, handler usecase.JobHandler) error {
	log.Printf("[Subscriber] subscribing to job type: %s", jobType)
	ch := make(chan []byte, 100) // バッファを持たせておく

	s.hub.mu.Lock()
	s.hub.channels[string(jobType)] = append(s.hub.channels[string(jobType)], ch)
	s.hub.mu.Unlock()

	// メッセージを待ち受けるループ
	go func() {
		for {
			select {
			case data := <-ch:
				log.Printf("[Subscriber] received job: type=%s", jobType)
				if err := handler(ctx, data); err != nil {
					log.Printf("[Subscriber] error handling job: type=%s err=%v", jobType, err)
				}
			case <-ctx.Done():
				log.Printf("[Subscriber] context cancelled for job type: %s", jobType)
				return
			}
		}
	}()

	return nil
}
