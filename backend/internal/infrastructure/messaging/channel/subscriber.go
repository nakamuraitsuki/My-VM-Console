package channel

import (
	"context"

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
	ch := make(chan []byte, 100) // バッファを持たせておく

	s.hub.mu.Lock()
	s.hub.channels[string(jobType)] = append(s.hub.channels[string(jobType)], ch)
	s.hub.mu.Unlock()

	// メッセージを待ち受けるループ
	go func() {
		for {
			select {
			case data := <-ch:
				if err := handler(ctx, data); err != nil {
					// エラーハンドリング（ログ出力など）
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	return nil
}
