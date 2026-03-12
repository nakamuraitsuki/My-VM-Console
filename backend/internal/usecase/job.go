package usecase

import "context"

type JobType string

const (
	JobTypeCreateInstance            JobType = "instance.create"
	JobTypeStopInstance              JobType = "instance.stop"
	JobTypeDeleteInstance            JobType = "instance.delete"
	JobTypeCreateVPCAndDefaultSubnet JobType = "vpc.create_with_default_subnet"
)

type JobPublisher interface {
	Publish(ctx context.Context, jobType JobType, payload any) error
}

type JobHandler func(ctx context.Context, payload any) error

type JobSubscriber interface {
	Subscribe(ctx context.Context, jobType JobType, handler JobHandler) error
}
