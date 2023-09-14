package db

import (
	"context"
	"wb_project/pkg/client/postgresql"
	"wb_project/pkg/logging"
	"wb_project/pkg/user"

	"github.com/jackc/pgconn"
)

type db struct {
	client postgresql.Client
	logger *logging.Logger
}

// Create implements user.Repository.
func (d *db) Create(ctx context.Context, u *user.User) error {
	q := `
	INSERT INTO public.users (track_number, order_uuid, entry, locale,
	 internal_signature, customer_id, delivery_service, shardkey,
	  sm_id, date_created, oof_shard)
	 VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	 RETURNING id
	 `
	if err := d.client.QueryRow(ctx, q, u.TrackNumber, u.OrderUid, u.Entry, u.Locale, u.InternalSignature, u.CustomerID,
		u.DeliveryService, u.Shardkey, u.SmId, u.DateCreated, u.OofShard).Scan(&u.ID); err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			d.logger.Errorf("SQL error: %s, Detail :%s, Where: %s, Code: %s, SQLState: %s", pgErr.Message, pgErr.Detail,
				pgErr.Where, pgErr.Code, pgErr.SQLState())
			return pgErr
		}
		return err
	}
	return nil
}

// FindAll implements user.Repository.
func (d *db) FindAll(ctx context.Context) (u []user.User, err error) {

	q := `
		SELECT * FROM public.users	
	`
	rows, err := d.client.Query(ctx, q)
	if err != nil {
		return nil, err
	}
	users := make([]user.User, 0)
	for rows.Next() {
		var use user.User
		err := rows.Scan(&use.ID, &use.OrderUid, &use.TrackNumber, &use.Entry, &use.Locale, &use.InternalSignature, &use.CustomerID,
			&use.DeliveryService, &use.Shardkey, &use.SmId, &use.DateCreated, &use.OofShard)
		if err != nil {
			return nil, err
		}
		users = append(users, use)

	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return users, nil
}

// FindOne implements user.Repository.
func (d *db) FindOne(ctx context.Context, id string) (user.User, error) {
	q := `
		SELECT * FROM public.users WHERE order_uuid = $1	
	`
	var use user.User
	err := d.client.QueryRow(ctx, q, id).Scan(&use.ID, &use.OrderUid, &use.TrackNumber, &use.Entry, &use.Locale,
		&use.InternalSignature, &use.CustomerID, &use.DeliveryService, &use.Shardkey, &use.SmId, &use.DateCreated, &use.OofShard)
	if err != nil {
		return user.User{}, err
	}
	return use, nil
}

// Update implements user.Repository.
func (d *db) Update(ctx context.Context, u user.User) error {
	q := `UPDATE public.users
	SET id = $1, 
	order_uuid = $2,
	track_number = $3,
	entry = $4,
	locale = $5,
	internal_signature = $6,
	customer_id = $7,
	delivery_service = $8,
	shardkey = $9,
	sm_id = $10,
	date_created = $11, 
	oof_shard = $12
	RETURNING id
	`
	if err := d.client.QueryRow(ctx, q, u.ID, u.OrderUid, u.TrackNumber, u.Entry, u.Locale, u.InternalSignature, u.CustomerID,
		u.DeliveryService, u.Shardkey, u.SmId, u.DateCreated, u.OofShard).Scan(&u.ID); err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			d.logger.Errorf("SQL error: %s, Detail :%s, Where: %s, Code: %s, SQLState: %s", pgErr.Message, pgErr.Detail,
				pgErr.Where, pgErr.Code, pgErr.SQLState())
			return pgErr
		}
		return err
	}
	return nil
}

// Delete implements user.Repository.
func (d *db) Delete(ctx context.Context, id string) error {
	return nil
}

func NewRepository(client postgresql.Client, logger *logging.Logger) user.Repository {
	return &db{
		client: client,
		logger: logger,
	}
}
