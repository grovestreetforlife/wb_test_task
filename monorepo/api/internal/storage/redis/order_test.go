package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redismock/v9"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
	"wb_test_task/api/internal/domain"
	"wb_test_task/libs/model"
)

func TestGetByID(t *testing.T) {
	client, mock := redismock.NewClientMock()
	cache := newOrderCache(client, 100)

	testCases := []struct {
		name      string
		mockInput struct {
			id string
		}
		mock           func(id string)
		expectedResult *model.Order
		wantErr        bool
		errMsg         string
	}{
		{
			name:      "OK",
			mockInput: struct{ id string }{id: "5d110e48-9e6b-4928-b436-14194b30d54f"},
			mock: func(id string) {
				mock.ExpectGet(fmt.Sprintf("%s:%s", orderObjectPrefix, id)).
					SetVal(`{"order_uid":"5d110e48-9e6b-4928-b436-14194b30d54f","track_number":"WBILMTESTTRACK3","entry":"WBIL","delivery":{"name":"Test Testov","phone":"+9720000000","zip":"2639809","city":"Kiryat Mozkin","address":"Ploshad Mira 15","region":"Kraiot","email":"test@gmail.com"},"payment":{"transaction":"5d110e48-9e6b-4928-b436-14194b30d54f","request_id":"5d110e48-9e6b-4928-b436-14194b30d54f","currency":"USD","provider":"wbpay","amount":1817,"payment_dt":1637907727,"bank":"alpha","delivery_cost":1500,"goods_total":317,"custom_fee":0},"items":[{"chrt_id":9934930,"track_number":"WBILMTESTTRACK3","price":453,"rid":"ab4219087a764ae0btest","name":"Mascaras","sale":30,"size":"0","total_price":317,"nm_id":2389212,"brand":"Vivienne Sabo","status":202}],"locale":"en","internal_signature":"","customer_id":"test","delivery_service":"meest","shard_key":"","sm_id":99,"date_created":"2021-11-26 06:22:19 +0000 UTC","oof_shard":"1"}`)
			},
			expectedResult: &model.Order{
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
			},
		},
		{
			name: "Order does not exists",
			mock: func(id string) {
				mock.ExpectGet(fmt.Sprintf("%s:%s", orderObjectPrefix, id)).RedisNil()
			},
			expectedResult: &model.Order{Items: []*model.Product{}},
			errMsg:         domain.ErrOrderNotExists.Error(),
			wantErr:        true,
		},
		{
			name: "Model and json fields mismatch",
			mock: func(id string) {
				mock.ExpectGet(fmt.Sprintf("%s:%s", orderObjectPrefix, id)).SetVal(``)
			},
			expectedResult: &model.Order{Items: []*model.Product{}},
			errMsg:         "unexpected end of JSON input",
			wantErr:        true,
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			test.mock(test.mockInput.id)

			order, err := cache.GetByID(context.Background(), test.mockInput.id)

			if test.wantErr {
				assert.EqualError(t, err, test.errMsg)
				assert.Equal(t, test.expectedResult, order)
			} else {
				assert.Equal(t, test.expectedResult, order)
				assert.NoError(t, err)
			}
		})
	}
}

func TestSet(t *testing.T) {
	client, mock := redismock.NewClientMock()
	cache := newOrderCache(client, 100)

	testCases := []struct {
		name      string
		mockInput struct {
			key   string
			order *model.Order
			ttl   int
		}
		mock    func(key string, order interface{}, ttl int)
		wantErr bool
		errMsg  string
	}{
		{
			name: "OK",
			mockInput: struct {
				key   string
				order *model.Order
				ttl   int
			}{
				key: "5d110e48-9e6b-4928-b436-14194b30d54f",
				order: &model.Order{
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
				},
				ttl: 100,
			},
			mock: func(key string, order interface{}, ttl int) {
				mock.ExpectSet(fmt.Sprintf("%s:%s", orderObjectPrefix, key), order, time.Duration(ttl)*time.Second).
					SetVal(`{"order_uid":"5d110e48-9e6b-4928-b436-14194b30d54f","track_number":"WBILMTESTTRACK3","entry":"WBIL","delivery":{"name":"Test Testov","phone":"+9720000000","zip":"2639809","city":"Kiryat Mozkin","address":"Ploshad Mira 15","region":"Kraiot","email":"test@gmail.com"},"payment":{"transaction":"5d110e48-9e6b-4928-b436-14194b30d54f","request_id":"5d110e48-9e6b-4928-b436-14194b30d54f","currency":"USD","provider":"wbpay","amount":1817,"payment_dt":1637907727,"bank":"alpha","delivery_cost":1500,"goods_total":317,"custom_fee":0},"items":[{"chrt_id":9934930,"track_number":"WBILMTESTTRACK3","price":453,"rid":"ab4219087a764ae0btest","name":"Mascaras","sale":30,"size":"0","total_price":317,"nm_id":2389212,"brand":"Vivienne Sabo","status":202}],"locale":"en","internal_signature":"","customer_id":"test","delivery_service":"meest","shard_key":"","sm_id":99,"date_created":"2021-11-26 06:22:19 +0000 UTC","oof_shard":"1"}`)
			},
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			data, err := json.Marshal(test.mockInput.order)
			if err != nil {
				t.Error(err)
			}

			test.mock(test.mockInput.key, data, test.mockInput.ttl)

			err = cache.Set(context.Background(), test.mockInput.key, test.mockInput.order)

			if test.wantErr {
				assert.EqualError(t, err, test.errMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
