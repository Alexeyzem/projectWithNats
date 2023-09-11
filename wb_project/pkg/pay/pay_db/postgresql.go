package paydb

import (
	"context"
	"wb_project/pkg/client/postgresql"
	"wb_project/pkg/logging"
	"wb_project/pkg/pay"

	"github.com/jackc/pgconn"
)

type db struct {
	client postgresql.Client
	logger *logging.Logger
}

// Create implements pay.Repository.
func (d *db) Create(ctx context.Context, p *pay.Pay) error {
	q := `
	INSERT INTO public.payment (user_id, transaction, request_id, currency,
		provider, amount, payment_dt, bank,
		delivery_cost, goods_total, custom_fee)
	 VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	 RETURNING id
	 `
	if err := d.client.QueryRow(ctx, q, p.UserId, p.Transaction, p.RequestId, p.Currency, p.Provider, p.Amount,
		p.PaymentDt, p.Bank, p.DeliveryCost, p.GoodsTotal, p.CustomFee).Scan(&p.ID); err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			d.logger.Errorf("SQL error: %s, Detail :%s, Where: %s, Code: %s, SQLState: %s", pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState())
			return pgErr
		}
		return err
	}
	return nil
}

// Delete implements pay.Repository.
func (d *db) Delete(ctx context.Context, id string) error {
	panic("unimplemented")
}

// FindAll implements pay.Repository.
func (d *db) FindAll(ctx context.Context) (p []pay.Pay, err error) {
	q := `
		SELECT * FROM public.payment	
	`
	rows, err := d.client.Query(ctx, q)
	if err != nil {
		return nil, err
	}
	pays := make([]pay.Pay, 0)
	for rows.Next() {
		var ps pay.Pay
		err := rows.Scan(&ps.ID, &ps.UserId, &ps.Transaction, &ps.RequestId, &ps.Currency, &ps.Provider,
			&ps.Amount, &ps.PaymentDt, &ps.Bank, &ps.DeliveryCost, &ps.GoodsTotal, &ps.CustomFee)
		if err != nil {
			return nil, err
		}
		pays = append(pays, ps)

	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return pays, nil
}

// FindOne implements pay.Repository.
func (d *db) FindOne(ctx context.Context, id string) (pay.Pay, error) {
	q := `
		SELECT * FROM public.payment WHERE user_id = &1	
	`
	var ps pay.Pay
	err := d.client.QueryRow(ctx, q, id).Scan(&ps.ID, &ps.UserId, &ps.Transaction, &ps.RequestId, &ps.Currency, &ps.Provider,
		&ps.Amount, &ps.PaymentDt, &ps.Bank, &ps.DeliveryCost, &ps.GoodsTotal, &ps.CustomFee)
	if err != nil {
		return pay.Pay{}, err
	}
	return ps, nil
}

// Update implements pay.Repository.
func (d *db) Update(ctx context.Context, p pay.Pay) error {
	panic("unimplemented")
}

func NewRepository(client postgresql.Client, logger *logging.Logger) pay.Repository {
	return &db{
		client: client,
		logger: logger,
	}
}
