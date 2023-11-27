package v1api

import (
	"context"
	"fmt"
	"github.com/dany-ykl/logger"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"wb_test_task/api/internal/common"
	mock_v1api "wb_test_task/api/internal/delivery/http/v1api/mocks"
	"wb_test_task/api/internal/domain"
	"wb_test_task/libs/model"
)

func init() {
	if err := logger.InitLogger(logger.Config{
		Namespace:   "test.order.api",
		Development: false,
		Filepath:    "",
		Level:       logger.InfoLevel,
	}); err != nil {
		log.Fatalln(err)
	}
}

func TestGetOrder(t *testing.T) {
	testCases := []struct {
		name                 string
		expectedStatusCode   int
		expectedResponseBody string
		mockInput            struct {
			orderID string
		}
		httpParam struct {
			key, value string
		}
		mockBehavior func(s *mock_v1api.MockorderService, ctx context.Context, id string)
	}{
		{
			name:               "OK",
			expectedStatusCode: http.StatusOK,
			mockInput:          struct{ orderID string }{orderID: "5d110e48-9e6b-4928-b436-14194b30d54f"},
			httpParam:          struct{ key, value string }{key: "id", value: "5d110e48-9e6b-4928-b436-14194b30d54f"},
			mockBehavior: func(s *mock_v1api.MockorderService, ctx context.Context, id string) {
				s.EXPECT().GetByID(ctx, id).Return(&model.Order{
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
				}, nil)
			},
			expectedResponseBody: fmt.Sprintf(`{"code":"OK","status":"ok","body":{"order_uid":"5d110e48-9e6b-4928-b436-14194b30d54f","track_number":"WBILMTESTTRACK3","entry":"WBIL","delivery":{"name":"Test Testov","phone":"+9720000000","zip":"2639809","city":"Kiryat Mozkin","address":"Ploshad Mira 15","region":"Kraiot","email":"test@gmail.com"},"payment":{"transaction":"5d110e48-9e6b-4928-b436-14194b30d54f","request_id":"5d110e48-9e6b-4928-b436-14194b30d54f","currency":"USD","provider":"wbpay","amount":1817,"payment_dt":1637907727,"bank":"alpha","delivery_cost":1500,"goods_total":317,"custom_fee":0},"items":[{"chrt_id":9934930,"track_number":"WBILMTESTTRACK3","price":453,"rid":"ab4219087a764ae0btest","name":"Mascaras","sale":30,"size":"0","total_price":317,"nm_id":2389212,"brand":"Vivienne Sabo","status":202}],"locale":"en","internal_signature":"","customer_id":"test","delivery_service":"meest","shard_key":"","sm_id":99,"date_created":"2021-11-26 06:22:19 +0000 UTC","oof_shard":"1"},"error":""}%s`, "\n"),
		},
		{
			name:               "Invalid uuid id",
			expectedStatusCode: http.StatusBadRequest,
			httpParam:          struct{ key, value string }{key: "id", value: "invalid"},
			mockInput:          struct{ orderID string }{orderID: "invalid"},
			mockBehavior: func(s *mock_v1api.MockorderService, ctx context.Context, id string) {
				s.EXPECT().GetByID(ctx, id).AnyTimes()
			},
			expectedResponseBody: fmt.Sprintf(`{"code":"Bad Request","status":"fail","body":null,"error":"invalid order uuid id"}%s`, "\n"),
		},
		{
			name:               "Null id http param",
			expectedStatusCode: http.StatusBadRequest,
			httpParam:          struct{ key, value string }{key: "", value: ""},
			mockInput:          struct{ orderID string }{orderID: ""},
			mockBehavior: func(s *mock_v1api.MockorderService, ctx context.Context, id string) {
				s.EXPECT().GetByID(ctx, id).AnyTimes()
			},
			expectedResponseBody: fmt.Sprintf(`{"code":"Bad Request","status":"fail","body":null,"error":"invalid order uuid id"}%s`, "\n"),
		},
		{
			name:               "Order does not exists",
			expectedStatusCode: http.StatusNotFound,
			mockInput:          struct{ orderID string }{orderID: "5d110e48-9e6b-4928-b436-14194b30d54f"},
			httpParam:          struct{ key, value string }{key: "id", value: "5d110e48-9e6b-4928-b436-14194b30d54f"},
			mockBehavior: func(s *mock_v1api.MockorderService, ctx context.Context, id string) {
				s.EXPECT().GetByID(ctx, id).Return(
					&model.Order{Items: []*model.Product{}},
					common.WrapError{Err: domain.ErrOrderNotExists, Msg: domain.ErrOrderNotExists.Error()},
				)
			},
			expectedResponseBody: fmt.Sprintf(`{"code":"Not Found","status":"fail","body":null,"error":"order does not exists"}%s`, "\n"),
		},
		{
			name:               "Interval Server Error",
			expectedStatusCode: http.StatusInternalServerError,
			mockInput:          struct{ orderID string }{orderID: "5d110e48-9e6b-4928-b436-14194b30d54f"},
			httpParam:          struct{ key, value string }{key: "id", value: "5d110e48-9e6b-4928-b436-14194b30d54f"},
			mockBehavior: func(s *mock_v1api.MockorderService, ctx context.Context, id string) {
				s.EXPECT().GetByID(ctx, id).Return(&model.Order{Items: []*model.Product{}}, errors.New("unexpected error"))
			},
			expectedResponseBody: fmt.Sprintf(`{"code":"Internal Server Error","status":"fail","body":null,"error":"unexpected error"}%s`, "\n"),
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			ct := gomock.NewController(t)
			defer ct.Finish()

			orderService := mock_v1api.NewMockorderService(ct)
			test.mockBehavior(orderService, context.Background(), test.mockInput.orderID)

			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rec := httptest.NewRecorder()

			e := echo.New()
			c := e.NewContext(req, rec)
			c.SetParamNames(test.httpParam.key)
			c.SetParamValues(test.httpParam.value)

			api := API{orderService: orderService}

			if assert.NoError(t, api.getOrder(c)) {
				assert.Equal(t, test.expectedStatusCode, rec.Code)
				assert.Equal(t, test.expectedResponseBody, rec.Body.String())
			}
		})
	}
}
