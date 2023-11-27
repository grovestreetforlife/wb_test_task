package services

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"testing"
	"wb_test_task/consumer/internal/domain"
	mock_services "wb_test_task/consumer/internal/services/mocks"
	"wb_test_task/libs/model"
)

func TestCreate(t *testing.T) {
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

	testCases := []struct {
		name      string
		mockInput struct {
			request *domain.OrderCreateRequest
		}
		mock func(storage *mock_services.MockorderStorage,
			cache *mock_services.MockorderCache, ctx context.Context,
			request *domain.OrderCreateRequest, order *model.Order)
		expectedResult *model.Order
		wantErr        bool
		errMsg         string
	}{
		{
			name:      "OK",
			mockInput: struct{ request *domain.OrderCreateRequest }{request: createOrderRequest},
			mock: func(storage *mock_services.MockorderStorage,
				cache *mock_services.MockorderCache, ctx context.Context,
				request *domain.OrderCreateRequest, order *model.Order) {
				storage.EXPECT().Create(ctx, request).Return(order, nil)
				cache.EXPECT().Set(ctx, request.OrderUid, order).Return(nil)
			},
			expectedResult: order,
			wantErr:        false,
			errMsg:         "",
		},
		{
			name:      "Error from storage",
			mockInput: struct{ request *domain.OrderCreateRequest }{request: createOrderRequest},
			mock: func(storage *mock_services.MockorderStorage,
				cache *mock_services.MockorderCache, ctx context.Context,
				request *domain.OrderCreateRequest, order *model.Order) {
				storage.EXPECT().Create(ctx, request).Return(&model.Order{Items: []*model.Product{}}, errors.New("error"))
			},
			expectedResult: &model.Order{Items: []*model.Product{}},
			wantErr:        true,
			errMsg:         "error",
		},
		{
			name:      "Error from cache",
			mockInput: struct{ request *domain.OrderCreateRequest }{request: createOrderRequest},
			mock: func(storage *mock_services.MockorderStorage,
				cache *mock_services.MockorderCache, ctx context.Context,
				request *domain.OrderCreateRequest, order *model.Order) {
				storage.EXPECT().Create(ctx, request).Return(order, nil)
				cache.EXPECT().Set(ctx, request.OrderUid, order).Return(errors.New("error"))
			},
			expectedResult: &model.Order{Items: []*model.Product{}},
			wantErr:        true,
			errMsg:         "error",
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			ct := gomock.NewController(t)
			defer ct.Finish()

			storage := mock_services.NewMockorderStorage(ct)
			cache := mock_services.NewMockorderCache(ct)

			test.mock(storage, cache, context.Background(), createOrderRequest, order)

			service := newOrderService(storage, cache)
			result, err := service.Create(context.Background(), createOrderRequest)

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
