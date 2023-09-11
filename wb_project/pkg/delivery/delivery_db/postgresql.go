package deliverydb

import (
	"context"
	"wb_project/pkg/client/postgresql"
	"wb_project/pkg/delivery"
	"wb_project/pkg/logging"

	"github.com/jackc/pgconn"
)

type db struct {
	client postgresql.Client
	logger *logging.Logger
}

// Create implements delivery.Repository.
func (d *db) Create(ctx context.Context, dv *delivery.Delivery) error {
	q := `
	INSERT INTO public.delivery (user_id, name, phone, zip,
		city, address, region, email)
	 VALUES($1, $2, $3, $4, $5, $6, $7, $8)
	 RETURNING id
	 `
	if err := d.client.QueryRow(ctx, q, dv.UserId, dv.Name, dv.Phone, dv.Zip,
		dv.City, dv.Address, dv.Region, dv.Email).Scan(&dv.ID); err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			d.logger.Errorf("SQL error: %s, Detail :%s, Where: %s, Code: %s, SQLState: %s", pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState())
			return pgErr
		}
		return err
	}
	return nil
}

// Delete implements delivery.Repository.
func (d *db) Delete(ctx context.Context, id string) error {
	panic("unimplemented")
}

// FindAll implements delivery.Repository.
func (d *db) FindAll(ctx context.Context) (dv []delivery.Delivery, err error) {
	q := `
		SELECT * FROM public.delivery	
	`
	rows, err := d.client.Query(ctx, q)
	if err != nil {
		return nil, err
	}
	dv = make([]delivery.Delivery, 0)
	for rows.Next() {
		var deliv delivery.Delivery
		err := rows.Scan(&deliv.ID, &deliv.UserId, &deliv.Name, &deliv.Phone, &deliv.Zip,
			&deliv.City, &deliv.Address, &deliv.Region, &deliv.Email)
		if err != nil {
			return nil, err
		}
		dv = append(dv, deliv)

	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return dv, nil
}

// FindOne implements delivery.Repository.
func (d *db) FindOne(ctx context.Context, id string) (delivery.Delivery, error) {
	q := `
		SELECT * FROM public.delivery WHERE user_id = &1	
	`
	var deliv delivery.Delivery
	err := d.client.QueryRow(ctx, q, id).Scan(&deliv.ID, &deliv.UserId, &deliv.Name, &deliv.Phone, &deliv.Zip,
		&deliv.City, &deliv.Address, &deliv.Region, &deliv.Email)
	if err != nil {
		return delivery.Delivery{}, err
	}
	return deliv, nil
}

// Update implements delivery.Repository.
func (d *db) Update(ctx context.Context, dv delivery.Delivery) error {
	panic("unimplemented")
}

func NewRepository(client postgresql.Client, logger *logging.Logger) delivery.Repository {
	return &db{
		client: client,
		logger: logger,
	}
}
