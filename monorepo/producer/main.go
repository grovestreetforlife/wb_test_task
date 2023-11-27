package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"log"
	"wb_test_task/libs/model"
)

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

var orderCreateRequestJSON = `	{
  "order_uid": "b563feb7b2b84b6test",
  "track_number": "WBILMTESTTRACK",
  "entry": "WBIL",
  "delivery": {
    "name": "Test Testov",
    "phone": "+9720000000",
    "zip": "2639809",
    "city": "Kiryat Mozkin",
    "address": "Ploshad Mira 15",
    "region": "Kraiot",
    "email": "test@gmail.com"
  },

  "payment": {
    "transaction": "b563feb7b2b84b6test",
    "request_id": "",
    "currency": "USD",
    "provider": "wbpay",
    "amount": 1817,
    "payment_dt": 1637907727,
    "bank": "alpha",
    "delivery_cost": 1500,
    "goods_total": 317,
    "custom_fee": 0
  },

  "items": [
    {
      "chrt_id": 9934930,
      "track_number": "WBILMTESTTRACK",
      "price": 453,
      "rid": "ab4219087a764ae0btest",
      "name": "Mascaras",
      "sale": 30,
      "size": "0",
      "total_price": 317,
      "nm_id": 2389212,
      "brand": "Vivienne Sabo",
      "status": 202
    },
    {
      "chrt_id": 9934931,
      "track_number": "WBILMTESTTRACK",
      "price": 453,
      "rid": "ab4219087a764ae0btest",
      "name": "Mascaras",
      "sale": 30,
      "size": "0",
      "total_price": 317,
      "nm_id": 2389212,
      "brand": "Vivienne Sabo",
      "status": 202
    }
  ],

  "locale": "en",
  "internal_signature": "",
  "customer_id": "test",
  "delivery_service": "meest",
  "shardkey": "9",
  "sm_id": 99,
  "date_created": "2021-11-26T06:22:19Z",
  "oof_shard": "1"
}`

const count = 100

func main() {
	producer, err := NewProducer(Config{Url: "http://localhost:4222"})
	if err != nil {
		log.Fatalln(err)
	}

	for i := 0; i < count; i++ {
		order := generateOrderCreateRequest(i)
		data, err := json.Marshal(order)
		if err != nil {
			log.Fatalln(err)
		}

		if err := producer.Publish(context.Background(), "order.create", data); err != nil {
			log.Fatalln(err)
		}
	}
}

func generateOrderCreateRequest(n int) OrderCreateRequest {
	var orderCreateRequest OrderCreateRequest
	if err := json.Unmarshal([]byte(orderCreateRequestJSON), &orderCreateRequest); err != nil {
		log.Fatalln(err)
	}

	orderUid := uuid.New().String()
	orderID := orderUid
	orderCreateRequest.OrderUid = orderUid
	orderCreateRequest.Payment.Transaction = orderID
	orderCreateRequest.Payment.RequestID = orderUid

	trackNumber := fmt.Sprintf("WBILMTESTTRACK%d", n)
	orderCreateRequest.TrackNumber = trackNumber
	for _, item := range orderCreateRequest.Items {
		item.TrackNumber = trackNumber
	}

	return orderCreateRequest
}
