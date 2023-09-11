package delivery

import "context"

type Repository interface {
	Create(ctx context.Context, dv *Delivery) error
	FindAll(ctx context.Context) (dv []Delivery, err error)
	FindOne(ctx context.Context, id string) (Delivery, error)
	Update(ctx context.Context, dv Delivery) error
	Delete(ctx context.Context, id string) error
}
