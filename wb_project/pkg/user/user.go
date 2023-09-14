package user

import (
	"wb_project/pkg/delivery"
	"wb_project/pkg/items"
	"wb_project/pkg/pay"
)

type User struct {
	ID                int               `json:"id"`
	OrderUid          string            `json:"order_uid"`
	TrackNumber       string            `json:"track_number"`
	Entry             string            `json:"entry"`
	Deliv             delivery.Delivery `json:"delivery"`
	Payment           pay.Pay           `json:"payment"`
	Items             []items.Item      `json:"items"`
	Locale            string            `json:"locale"`
	InternalSignature string            `json:"internal_signature"`
	CustomerID        string            `json:"customer_id"`
	DeliveryService   string            `json:"delivery_service"`
	Shardkey          string            `json:"shardkey"`
	SmId              int               `json:"sm_id"`
	DateCreated       string            `json:"date_created"`
	OofShard          string            `json:"oof_shard"`
}
