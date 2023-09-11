package pay

import "context"

type Repository interface {
	Create(ctx context.Context, p *Pay) error
	FindAll(ctx context.Context) (p []Pay, err error)
	FindOne(ctx context.Context, id string) (Pay, error)
	Update(ctx context.Context, p Pay) error
	Delete(ctx context.Context, id string) error
}
