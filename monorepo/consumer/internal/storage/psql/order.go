package psql

import (
	"context"
	"fmt"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
	"wb_test_task/consumer/internal/common"
	"wb_test_task/consumer/internal/domain"
	"wb_test_task/libs/model"
)

type orderStorage struct {
	pool pool
}

func newOrderStorage(pool pool) *orderStorage {
	return &orderStorage{pool: pool}
}

// Create создание заказа
func (o *orderStorage) Create(ctx context.Context, request *domain.OrderCreateRequest) (*model.Order, error) {
	tx, err := o.pool.Begin(ctx)
	if err != nil {
		return &model.Order{}, common.WrapError{Err: err, Msg: "fail to create transaction"}
	}
	defer tx.Rollback(context.Background())

	queryOrderCreate := `
		INSERT INTO orders(order_uid, track_number, entry, locale, internal_signature,
		                   customer_id, delivery_service, shardkey, sm_id, oof_shard,
		                   date_created)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`

	_, err = tx.Exec(ctx, queryOrderCreate, request.OrderUid, request.TrackNumber, request.Entry, request.Locale,
		request.InternalSignature, request.CustomerID, request.DeliveryService, request.ShardKey, request.SmID,
		request.OofShard, request.DateCreated)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case domain.CodeErrDuplicateKey:
				return &model.Order{}, common.WrapError{Err: domain.ErrOrderAlreadyExists,
					Msg: fmt.Sprintf("hint: %s", pgErr.Detail),
				}
			case domain.CodeErrConstraintLenValue:
				return &model.Order{}, common.WrapError{Err: domain.ErrInvalidOrderValue, Msg: pgErr.Message}
			case domain.CodeErrInvalidSyntax:
				return &model.Order{}, common.WrapError{Err: domain.ErrInvalidValue, Msg: pgErr.Message}
			case domain.CodeErrForeignKey:
				return &model.Order{}, common.WrapError{Err: domain.ErrUserDoesNotExists, Msg: domain.ErrUserDoesNotExists.Error()}
			}
		}

		return &model.Order{}, common.WrapError{Err: err, Msg: "fail to create order"}
	}

	// создание транзакции
	if err := o.createPayment(ctx, tx, request); err != nil {
		return &model.Order{}, err
	}

	// создание доставки
	if err := o.createDelivery(ctx, tx, request); err != nil {
		return &model.Order{}, err
	}

	// создание продуктов
	if err := o.createProducts(ctx, tx, request); err != nil {
		return &model.Order{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		return &model.Order{}, common.WrapError{Err: err, Msg: "fail to commit transaction"}
	}

	return mapOrderCreateRequestToOrder(request), nil
}

// createPayment создание транзакции
func (o *orderStorage) createPayment(ctx context.Context, tx pgx.Tx, request *domain.OrderCreateRequest) error {
	queryOrderPaymentCreate := `
		INSERT INTO transaction(id, request_id, currency, provider, amount, payment_dt,
		                        bank, delivery_cost, goods_total, custom_fee)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`
	_, err := tx.Exec(ctx, queryOrderPaymentCreate, request.OrderUid, request.Payment.RequestID,
		request.Payment.Currency, request.Payment.Provider, request.Payment.Amount, request.Payment.PaymentDt,
		request.Payment.Bank, request.Payment.DeliveryCost, request.Payment.GoodsTotal, request.Payment.CustomFee)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case domain.CodeErrDuplicateKey:
				return common.WrapError{Err: domain.ErrOrderAlreadyExists,
					Msg: fmt.Sprintf("hint: %s", pgErr.Detail)}
			case domain.CodeErrConstraintLenValue:
				return common.WrapError{Err: domain.ErrInvalidOrderValue, Msg: pgErr.Message}
			case domain.CodeErrInvalidSyntax:
				return common.WrapError{Err: domain.ErrInvalidValue,
					Msg: fmt.Sprintf("hint: %s", pgErr.Message)}
			case domain.CodeErrForeignKey:
				return common.WrapError{Err: domain.ErrOrderDoesNotExists, Msg: domain.ErrOrderDoesNotExists.Error()}
			}
		}

		return common.WrapError{Err: err, Msg: "fail to create order"}
	}

	return nil
}

// createDelivery создание доставки
func (o *orderStorage) createDelivery(ctx context.Context, tx pgx.Tx, request *domain.OrderCreateRequest) error {
	queryOrderDeliveryCreate := `
		INSERT INTO delivery(order_uid, name, phone, zip, city, address, region, email)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	_, err := tx.Exec(ctx, queryOrderDeliveryCreate, request.OrderUid, request.Delivery.Name, request.Delivery.Phone,
		request.Delivery.Zip, request.Delivery.City, request.Delivery.Address, request.Delivery.Region, request.Delivery.Email)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case domain.CodeErrDuplicateKey:
				return common.WrapError{Err: domain.ErrOrderAlreadyExists, Msg: domain.ErrOrderAlreadyExists.Error()}
			case domain.CodeErrConstraintLenValue:
				return common.WrapError{Err: domain.ErrInvalidOrderValue, Msg: pgErr.Message}
			case domain.CodeErrInvalidSyntax:
				return common.WrapError{Err: domain.ErrInvalidValue, Msg: pgErr.Message}
			case domain.CodeErrForeignKey:
				return common.WrapError{Err: domain.ErrOrderDoesNotExists, Msg: domain.ErrOrderDoesNotExists.Error()}
			}
		}

		return common.WrapError{Err: err, Msg: "fail to create order"}
	}

	return nil
}

// createProducts создание продукта
func (o *orderStorage) createProducts(ctx context.Context, tx pgx.Tx, request *domain.OrderCreateRequest) error {
	rows := make([][]interface{}, 0, len(request.Items))

	for _, product := range request.Items {
		rows = append(rows, []interface{}{product.ChrtID, product.TrackNumber, product.Price, product.Rid,
			product.Name, product.Sale, product.Size, product.TotalPrice, product.NmID, product.Brand, product.Status})
	}

	_, err := tx.CopyFrom(
		ctx,
		pgx.Identifier{"product"},
		[]string{"chrt_id", "track_number", "price", "rid", "name", "sale", "size",
			"total_price", "nm_id", "brand", "status"},
		pgx.CopyFromRows(rows),
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case domain.CodeErrDuplicateKey:
				return common.WrapError{Err: domain.ErrOrderAlreadyExists, Msg: domain.ErrOrderAlreadyExists.Error()}
			case domain.CodeErrConstraintLenValue:
				return common.WrapError{Err: domain.ErrInvalidOrderValue, Msg: pgErr.Message}
			case domain.CodeErrInvalidSyntax:
				return common.WrapError{Err: domain.ErrInvalidValue, Msg: pgErr.Message}
			case domain.CodeErrForeignKey:
				return common.WrapError{Err: domain.ErrOrderDoesNotExists, Msg: domain.ErrOrderDoesNotExists.Error()}
			}
		}

		return common.WrapError{Err: err, Msg: "fail to create order"}
	}

	return nil
}

func mapOrderCreateRequestToOrder(request *domain.OrderCreateRequest) *model.Order {
	return &model.Order{
		OrderUid:          request.OrderUid,
		TrackNumber:       request.TrackNumber,
		Entry:             request.Entry,
		Delivery:          request.Delivery,
		Payment:           request.Payment,
		Items:             request.Items,
		Locale:            request.Locale,
		InternalSignature: request.InternalSignature,
		CustomerID:        request.CustomerID,
		DeliveryService:   request.DeliveryService,
		ShardKey:          request.ShardKey,
		SmID:              request.SmID,
		DateCreated:       request.DateCreated,
		OofShard:          request.OofShard,
	}
}
