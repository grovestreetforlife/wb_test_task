package consumer

import (
	"context"
	"encoding/json"
	"github.com/dany-ykl/logger"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
	"wb_test_task/consumer/internal/config"
	mock_consumer "wb_test_task/consumer/internal/consumer/mocks"
	"wb_test_task/consumer/internal/domain"
	"wb_test_task/libs/model"
)

func init() {
	if err := logger.InitLogger(logger.Config{
		Namespace:   "test.consumer",
		Development: true,
		Level:       logger.InfoLevel,
	}); err != nil {
		log.Fatalln(err)
	}
}

func TestOrderCreateHandler(t *testing.T) {
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
	data, err := json.Marshal(createOrderRequest)
	if err != nil {
		t.Error(err)
	}

	testCases := []struct {
		name      string
		mockInput struct {
			msg     *Msg
			request *domain.OrderCreateRequest
		}
		mock    func(service *mock_consumer.MockorderService, ctx context.Context, request *domain.OrderCreateRequest)
		wantErr bool
		errMsg  string
	}{
		{
			name: "OK",
			mockInput: struct {
				msg     *Msg
				request *domain.OrderCreateRequest
			}{msg: &Msg{Subject: "order.create", Data: data}, request: createOrderRequest},
			mock: func(service *mock_consumer.MockorderService, ctx context.Context, request *domain.OrderCreateRequest) {
				service.EXPECT().Create(ctx, request).Return(order, nil)
			},
			wantErr: false,
			errMsg:  "",
		},
		{
			name: "Service error",
			mockInput: struct {
				msg     *Msg
				request *domain.OrderCreateRequest
			}{msg: &Msg{Subject: "order.create", Data: data}, request: createOrderRequest},
			mock: func(service *mock_consumer.MockorderService, ctx context.Context, request *domain.OrderCreateRequest) {
				service.EXPECT().Create(ctx, request).Return(&model.Order{Items: []*model.Product{}}, errors.New("error"))
			},
			wantErr: true,
			errMsg:  "fail to create order: error",
		},
		{
			name: "Unmarshal error",
			mockInput: struct {
				msg     *Msg
				request *domain.OrderCreateRequest
			}{msg: &Msg{Subject: "order.create", Data: []byte(`error`)}, request: &domain.OrderCreateRequest{}},
			mock: func(service *mock_consumer.MockorderService, ctx context.Context, request *domain.OrderCreateRequest) {
			},
			wantErr: true,
			errMsg:  "fail to unmarshal msg: invalid character 'e' looking for beginning of value",
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			ct := gomock.NewController(t)
			defer ct.Finish()

			orderService := mock_consumer.NewMockorderService(ct)
			test.mock(orderService, context.Background(), createOrderRequest)

			consumer := Consumer{orderService: orderService}
			err := consumer.orderCreateHandler(context.Background(), test.mockInput.msg)

			if test.wantErr {
				assert.EqualError(t, err, test.errMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestNew(t *testing.T) {
	ct := gomock.NewController(t)
	defer ct.Finish()

	cfg := config.NatsConsumer{
		Url:                  "nats://127.0.0.1:4222",
		Subjects:             []string{"order.create"},
		RetryOfFailedConnect: true,
		StreamName:           "orders",
		CountConsumers:       2,
	}
	orderService := mock_consumer.NewMockorderService(ct)

	_, err := New(cfg, orderService)

	assert.NoError(t, err)
}

func TestStart(t *testing.T) {
	if err := logger.InitLogger(logger.Config{
		Namespace:   "test.consumer",
		Development: true,
		Level:       logger.InfoLevel,
	}); err != nil {
		t.Error(err)
	}

	ct := gomock.NewController(t)
	defer ct.Finish()

	cfg := config.NatsConsumer{
		Url:                  "nats://127.0.0.1:4222",
		Subjects:             []string{"order.create"},
		RetryOfFailedConnect: true,
		StreamName:           "orders",
		CountConsumers:       2,
	}
	orderService := mock_consumer.NewMockorderService(ct)

	consumer, err := New(cfg, orderService)
	assert.NoError(t, err)

	err = consumer.Start(context.Background())
	assert.NoError(t, err)
}
