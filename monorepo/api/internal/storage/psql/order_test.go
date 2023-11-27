package psql

import (
	"context"
	"github.com/pashagolub/pgxmock/v3"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
	"wb_test_task/libs/model"
)

func TestGetByID(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Error(err)
	}
	defer mock.Close()

	storage := newOrderStorage(mock)

	products := []*model.Product{
		{
			ChrtID:      9934930,
			TrackNumber: "WBILMTESTTRACK3",
			Price:       453,
			Rid:         "ab4219087a764ae0btest",
			Name:        "Mascaras",
			Sale:        30,
			Size:        "0",
			TotalPrice:  317,
			NmID:        2389212,
			Brand:       "Vivienne Sabo",
			Status:      202,
		},
	}
	order := &model.Order{
		OrderUid:    "5d110e48-9e6b-4928-b436-14194b30d54f",
		TrackNumber: "WBILMTESTTRACK3",
		Entry:       "WBIL",
		Delivery: model.Delivery{
			Name:    "Test Testov",
			Phone:   "+9720000000",
			Zip:     "2639809",
			City:    "Kiryat Mozkin",
			Address: "Ploshad Mira 15",
			Region:  "Kraiot",
			Email:   "test@gmail.com",
		},
		Payment: model.Payment{
			Transaction:  "5d110e48-9e6b-4928-b436-14194b30d54f",
			RequestID:    "5d110e48-9e6b-4928-b436-14194b30d54f",
			Currency:     "USD",
			Provider:     "wbpay",
			Amount:       1817,
			PaymentDt:    1637907727,
			Bank:         "alpha",
			DeliveryCost: 1500,
			GoodsTotal:   317,
			CustomFee:    0,
		},
		Items:             products,
		Locale:            "en",
		InternalSignature: "",
		CustomerID:        "test",
		DeliveryService:   "meest",
		ShardKey:          "",
		SmID:              99,
		DateCreated:       "2021-11-26 06:22:19 +0000 UTC",
		OofShard:          "1",
	}

	selectOrderByIDQuery := `
	SELECT o.order_uid, o.track_number, o.entry, o.locale, o.internal_signature,
       			o.customer_id, o.delivery_service, o.shardkey, o.sm_id, o.oof_shard,
       			o.date_created, d.name, d.phone, d.zip, d.city, d.address, d.region, d.email,
       			t.id, t.request_id, t.currency, t.provider, t.amount, t.payment_dt,
       			t.bank, t.delivery_cost, t.goods_total, t.custom_fee
				FROM orders o
				JOIN delivery d
    				ON o.order_uid=d.order_uid
				JOIN transaction t
    				ON o.order_uid=t.id
				WHERE o.order_uid
	`
	selectItemsByTrackNumberQuery := `
		
					SELECT chrt_id, track_number, price, rid, name, sale,
		   			 size, total_price, nm_id, brand, status
					FROM product
					WHERE track_number
					
	`

	orderRows := []string{
		"o.order_uid", "o.track_number", "o.entry", "o.locale", "o.internal_signature",
		"o.customer_id", "o.delivery_service", "o.shardkey", "o.sm_id", "o.oof_shard",
		"o.date_created", "d.name", "d.phone", "d.zip", "d.city", "d.address", "d.region", "d.email",
		"t.id", "t.request_id", "t.currency", "t.provider", "t.amount", "t.payment_dt",
		"t.bank", "t.delivery_cost", "t.goods_total", "t.custom_fee",
	}
	itemRows := []string{"chrt_id", "track_number", "price", "rid", "name", "sale",
		"size", "total_price", "nm_id", "brand", "status"}

	testCases := []struct {
		name      string
		mock      func(ctx context.Context, id, trackNumber string)
		mockInput struct {
			id          string
			trackNumber string
		}
		expectedResult *model.Order
		wantErr        bool
		errMsg         string
	}{
		{
			name: "OK",
			mock: func(ctx context.Context, id, trackNumber string) {
				// select order, payment, delivery
				rows := mock.NewRows(orderRows)
				dateCreated, err := time.Parse("2006-01-02 15:04:05 -0700 MST", "2021-11-26 06:22:19 +0000 UTC")
				if err != nil {
					t.Error(err)
				}
				rows.AddRow(
					"5d110e48-9e6b-4928-b436-14194b30d54f", "WBILMTESTTRACK3", "WBIL", "en", "", "test", "meest",
					"", 99, "1", dateCreated, "Test Testov", "+9720000000", "2639809", "Kiryat Mozkin",
					"Ploshad Mira 15", "Kraiot", "test@gmail.com", "5d110e48-9e6b-4928-b436-14194b30d54f", "5d110e48-9e6b-4928-b436-14194b30d54f",
					"USD", "wbpay", float64(1817), int64(1637907727), "alpha", float64(1500), 317, 0,
				)
				mock.ExpectQuery(selectOrderByIDQuery).WithArgs(id).WillReturnRows(rows)

				// select items
				itemRows := mock.NewRows(itemRows)
				itemRows.AddRow(int64(9934930), "WBILMTESTTRACK3", float64(453), "ab4219087a764ae0btest", "Mascaras", 30,
					"0", float64(317), int64(2389212), "Vivienne Sabo", 202)
				mock.ExpectQuery(selectItemsByTrackNumberQuery).WithArgs(trackNumber).WillReturnRows(itemRows)
			},
			mockInput: struct {
				id          string
				trackNumber string
			}{id: "5d110e48-9e6b-4928-b436-14194b30d54f", trackNumber: "WBILMTESTTRACK3"},
			expectedResult: order,
		},
		{
			name: "Order does not exists",
			mock: func(ctx context.Context, id, trackNumber string) {
				// select order, payment, delivery
				rows := mock.NewRows(orderRows)
				mock.ExpectQuery(selectOrderByIDQuery).WithArgs(id).WillReturnRows(rows)
			},
			mockInput: struct {
				id          string
				trackNumber string
			}{id: "5d110e48-9e6b-4928-b436-14194b30d54f", trackNumber: "WBILMTESTTRACK3"},
			wantErr:        true,
			expectedResult: &model.Order{Items: []*model.Product{}},
			errMsg:         "order does not exists",
		},
		{
			name: "Unexpected error",
			mock: func(ctx context.Context, id, trackNumber string) {
				// select order, payment, delivery
				mock.ExpectQuery(selectOrderByIDQuery).WithArgs(id).WillReturnError(errors.New("error"))
			},
			mockInput: struct {
				id          string
				trackNumber string
			}{id: "5d110e48-9e6b-4928-b436-14194b30d54f", trackNumber: "WBILMTESTTRACK3"},
			wantErr:        true,
			expectedResult: &model.Order{Items: []*model.Product{}},
			errMsg:         "error",
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			test.mock(context.Background(), test.mockInput.id, test.mockInput.trackNumber)

			result, err := storage.GetByID(context.Background(), test.mockInput.id)

			if test.wantErr {
				assert.EqualError(t, err, test.errMsg)
				assert.Equal(t, test.expectedResult, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.expectedResult, result)
			}
		})
	}
}
