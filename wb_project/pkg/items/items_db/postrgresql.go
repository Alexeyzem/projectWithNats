package itemsdb

import (
	"context"
	"fmt"
	"wb_project/pkg/client/postgresql"
	"wb_project/pkg/items"
	"wb_project/pkg/logging"

	"github.com/jackc/pgconn"
)

type db struct {
	client postgresql.Client
	logger *logging.Logger
}

// Create implements items.Repository.
func (d *db) Create(ctx context.Context, i *items.Item) error {
	q := `
	INSERT INTO public.items (chrt_id, track_number, price, rid, name,
		sale, size, total_price, nm_id, brand, status)
	 VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	 RETURNING id
	 `
	if err := d.client.QueryRow(ctx, q, i.ChrtId, i.TrackNumber, i.Price, i.Rid, i.Name,
		i.Sale, i.Size, i.TotalPrice, i.NmId, i.Brand, i.Status).Scan(&i.ID); err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			d.logger.Errorf("SQL error: %s, Detail :%s, Where: %s, Code: %s, SQLState: %s", pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState())
			return pgErr
		}
		return err
	}
	return nil
}

// Delete implements items.Repository.
func (d *db) Delete(ctx context.Context, id string) error {
	panic("unimplemented")
}

// FindAll implements items.Repository.
func (d *db) FindAll(ctx context.Context) (i []items.Item, err error) {
	q := `
		SELECT * FROM public.items	
	`
	i, err = d.find(ctx, q)
	if err != nil {
		return nil, err
	}
	return i, nil
}

// FindAllOfOneUser implements items.Repository.
func (d *db) FindAllOfOneUser(ctx context.Context, track_number string) ([]items.Item, error) {
	q := `
		SELECT * FROM public.items WHERE track_number = $1
	`
	i, err := d.find(ctx, q, track_number)
	if err != nil {
		return nil, err
	}
	return i, nil
}

// Update implements items.Repository.
func (d *db) Update(ctx context.Context, i items.Item) error {
	panic("unimplemented")
}

func NewRepository(client postgresql.Client, logger *logging.Logger) items.Repository {
	return &db{
		client: client,
		logger: logger,
	}
}

func (d *db) find(ctx context.Context, q string, track_number ...interface{}) ([]items.Item, error) {
	rows, err := d.client.Query(ctx, q, track_number[0])
	if err != nil {
		d.logger.Error(err)
		return nil, err
	}
	i := make([]items.Item, 0)
	for rows.Next() {
		var it items.Item
		err := rows.Scan(&it.ID, &it.ChrtId, &it.TrackNumber, &it.Price, &it.Rid, &it.Name,
			&it.Sale, &it.Size, &it.TotalPrice, &it.NmId, &it.Brand, &it.Status)
		if err != nil {
			return nil, err
		}
		i = append(i, it)

	}
	fmt.Println(i)
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return i, nil
}
