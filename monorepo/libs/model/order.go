package model

type Order struct {
	OrderUid          string     `json:"order_uid"`
	TrackNumber       string     `json:"track_number"`
	Entry             string     `json:"entry"`
	Delivery          Delivery   `json:"delivery"`
	Payment           Payment    `json:"payment"`
	Items             []*Product `json:"items"`
	Locale            string     `json:"locale"`
	InternalSignature string     `json:"internal_signature"`
	CustomerID        string     `json:"customer_id"`
	DeliveryService   string     `json:"delivery_service"`
	ShardKey          string     `json:"shard_key"`
	SmID              int        `json:"sm_id"`
	DateCreated       string     `json:"date_created"`
	OofShard          string     `json:"oof_shard"`
}
