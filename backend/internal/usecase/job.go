package usecase

import (
	"context"
	"encoding/json"
)

type JobType string

const (
	JobTypeCreateInstance            JobType = "instance.create"
	JobTypeStopInstance              JobType = "instance.stop"
	JobTypeDeleteInstance            JobType = "instance.delete"
	JobTypeCreateVPCAndDefaultSubnet JobType = "vpc.create_with_default_subnet"
)

type JobPublisher interface {
	Publish(ctx context.Context, jobType JobType, payload []byte) error
}

type JobHandler func(ctx context.Context, payload []byte) error

type JobSubscriber interface {
	Subscribe(ctx context.Context, jobType JobType, handler JobHandler) error
}

func Bind[T any](
	ctx context.Context,
	s JobSubscriber,
	jobType JobType,
	fn func(context.Context, T) error,
) error {
	handler := func(ctx context.Context, payload []byte) error {
		var data T
		if err := json.Unmarshal(payload, &data); err != nil {
			return err
		}
		return fn(ctx, data)
	}
	return s.Subscribe(ctx, jobType, handler)
}
