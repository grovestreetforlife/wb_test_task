package psql

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/pashagolub/pgxmock/v3"
	"github.com/stretchr/testify/assert"
	"testing"
	"wb_test_task/consumer/internal/domain"
	"wb_test_task/libs/model"
)

func TestGetByID(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Error(err)
	}
	defer mock.Close()

	createOrderRequest := &domain.OrderCreateRequest{
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
		Items: []*model.Product{
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
		},
		Locale:            "en",
		InternalSignature: "",
		CustomerID:        "test",
		DeliveryService:   "meest",
		ShardKey:          "",
		SmID:              99,
		DateCreated:       "2021-11-26 06:22:19 +0000 UTC",
		OofShard:          "1",
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
		Items: []*model.Product{
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
		},
		Locale:            "en",
		InternalSignature: "",
		CustomerID:        "test",
		DeliveryService:   "meest",
		ShardKey:          "",
		SmID:              99,
		DateCreated:       "2021-11-26 06:22:19 +0000 UTC",
		OofShard:          "1",
	}

	storage := newOrderStorage(mock)

	createOrderQuery := `INSERT INTO orders`
	createTransactionQuery := `INSERT INTO transaction`
	createDeliveryQuery := `INSERT INTO delivery`

	productRows := []string{"chrt_id", "track_number", "price", "rid", "name", "sale", "size",
		"total_price", "nm_id", "brand", "status"}

	testCases := []struct {
		name           string
		mock           func()
		expectedResult *model.Order
		wantErr        bool
		errMsg         string
	}{
		{
			name: "OK",
			mock: func() {
				mock.ExpectBegin()
				mock.ExpectExec(createOrderQuery).WithArgs("5d110e48-9e6b-4928-b436-14194b30d54f", "WBILMTESTTRACK3", "WBIL", "en", "", "test", "meest",
					"", 99, "1", "2021-11-26 06:22:19 +0000 UTC").WillReturnResult(pgxmock.NewResult("INSERT", 1))

				mock.ExpectExec(createTransactionQuery).WithArgs("5d110e48-9e6b-4928-b436-14194b30d54f", "5d110e48-9e6b-4928-b436-14194b30d54f",
					"USD", "wbpay", float64(1817), int64(1637907727), "alpha", float64(1500), 317, 0).WillReturnResult(pgxmock.NewResult("INSERT", 1))

				mock.ExpectExec(createDeliveryQuery).WithArgs("5d110e48-9e6b-4928-b436-14194b30d54f", "Test Testov", "+9720000000", "2639809", "Kiryat Mozkin",
					"Ploshad Mira 15", "Kraiot", "test@gmail.com").WillReturnResult(pgxmock.NewResult("INSERT", 1))

				mock.ExpectCopyFrom(pgx.Identifier{"product"}, productRows).WillReturnResult(1)
				mock.ExpectCommit()
			},
			expectedResult: order,
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			test.mock()

			result, err := storage.Create(context.Background(), createOrderRequest)

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
