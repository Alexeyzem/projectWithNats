package items

import "context"

type Repository interface {
	Create(ctx context.Context, i *Item) error
	FindAll(ctx context.Context) (i []Item, err error)
	FindAllOfOneUser(ctx context.Context, id string) ([]Item, error)
	Update(ctx context.Context, i Item) error
	Delete(ctx context.Context, id string) error
}
