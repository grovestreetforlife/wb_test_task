package services

import (
	"context"
	"github.com/dany-ykl/logger"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"log"
	"runtime"
	"testing"
	"wb_test_task/api/internal/common"
	"wb_test_task/api/internal/domain"
	mock_services "wb_test_task/api/internal/services/mocks"
	"wb_test_task/libs/model"
)

func init() {
	if err := logger.InitLogger(logger.Config{
		Namespace:   "test.order.service",
		Development: false,
		Filepath:    "",
		Level:       logger.InfoLevel,
	}); err != nil {
		log.Fatalln(err)
	}
}

func TestGetByID(t *testing.T) {
	runtime.GOMAXPROCS(1)

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
		mock      func(cache *mock_services.MockorderCache, storage *mock_services.MockorderStorage, ctx context.Context, id string, order *model.Order)
		mockInput struct {
			id    string
			order *model.Order
		}
		expectedResult *model.Order
		wantErr        bool
		errMsg         string
	}{
		{
			name: "OK. Order exists in cache",
			mockInput: struct {
				id    string
				order *model.Order
			}{id: "5d110e48-9e6b-4928-b436-14194b30d54f"},
			mock: func(cache *mock_services.MockorderCache, storage *mock_services.MockorderStorage, ctx context.Context, id string, order *model.Order) {
				cache.EXPECT().GetByID(ctx, id).Return(order, nil)
			},
			expectedResult: order,
			wantErr:        false,
			errMsg:         "",
		},
		{
			name: "OK. Order does not exists in cache and storage",
			mockInput: struct {
				id    string
				order *model.Order
			}{id: "5d110e48-9e6b-4928-b436-14194b30d54f", order: order},
			mock: func(cache *mock_services.MockorderCache, storage *mock_services.MockorderStorage, ctx context.Context, id string, order *model.Order) {
				cache.EXPECT().GetByID(ctx, id).Return(&model.Order{Items: []*model.Product{}}, common.WrapError{Err: domain.ErrOrderNotExists, Msg: domain.ErrOrderNotExists.Error()})
				storage.EXPECT().GetByID(ctx, id).Return(&model.Order{Items: []*model.Product{}}, common.WrapError{Err: domain.ErrOrderNotExists, Msg: domain.ErrOrderNotExists.Error()})
			},
			expectedResult: &model.Order{Items: []*model.Product{}},
			wantErr:        true,
			errMsg:         domain.ErrOrderNotExists.Error(),
		},
		{
			name: "OK. Error from cache",
			mockInput: struct {
				id    string
				order *model.Order
			}{id: "5d110e48-9e6b-4928-b436-14194b30d54f", order: order},
			mock: func(cache *mock_services.MockorderCache, storage *mock_services.MockorderStorage, ctx context.Context, id string, order *model.Order) {
				cache.EXPECT().GetByID(ctx, id).Return(&model.Order{Items: []*model.Product{}}, errors.New("unexpected error"))
				storage.EXPECT().GetByID(ctx, id).Return(order, nil)
				cache.EXPECT().Set(ctx, id, order).Return(nil).AnyTimes()
			},
			expectedResult: order,
			wantErr:        false,
			errMsg:         "",
		},
		{
			name: "Order does not exists in cache but exist in storage",
			mockInput: struct {
				id    string
				order *model.Order
			}{id: "5d110e48-9e6b-4928-b436-14194b30d54f", order: order},
			mock: func(cache *mock_services.MockorderCache, storage *mock_services.MockorderStorage, ctx context.Context, id string, order *model.Order) {
				cache.EXPECT().GetByID(ctx, id).Return(&model.Order{Items: []*model.Product{}}, common.WrapError{Err: domain.ErrOrderNotExists, Msg: domain.ErrOrderNotExists.Error()})
				storage.EXPECT().GetByID(ctx, id).Return(order, nil)
			},
			expectedResult: order,
			wantErr:        false,
			errMsg:         "",
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			ct := gomock.NewController(t)
			defer ct.Finish()

			cache := mock_services.NewMockorderCache(ct)
			storage := mock_services.NewMockorderStorage(ct)

			test.mock(cache, storage, context.Background(), test.mockInput.id, order)

			service := newOrderService(storage, cache)
			result, err := service.GetByID(context.Background(), test.mockInput.id)

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
