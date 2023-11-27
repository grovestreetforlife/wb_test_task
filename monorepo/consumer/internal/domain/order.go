package domain

import "wb_test_task/libs/model"

type OrderCreateRequest struct {
	OrderUid          string           `json:"order_uid"`
	TrackNumber       string           `json:"track_number"`
	Entry             string           `json:"entry"`
	Delivery          model.Delivery   `json:"delivery"`
	Payment           model.Payment    `json:"payment"`
	Items             []*model.Product `json:"items"`
	Locale            string           `json:"locale"`
	InternalSignature string           `json:"internal_signature"`
	CustomerID        string           `json:"customer_id"`
	DeliveryService   string           `json:"delivery_service"`
	ShardKey          string           `json:"shard_key"`
	SmID              int              `json:"sm_id"`
	DateCreated       string           `json:"date_created"`
	OofShard          string           `json:"oof_shard"`
}
