package psql

import (
	"context"
	"github.com/dany-ykl/tracer"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/pkg/errors"
	"go.opentelemetry.io/otel/attribute"
	"time"
	"wb_test_task/api/internal/common"
	"wb_test_task/api/internal/domain"
	"wb_test_task/libs/model"
)

type orderStorage struct {
	pool pool
}

func newOrderStorage(pool pool) *orderStorage {
	return &orderStorage{pool: pool}
}

// GetByID вернуть заказ по его id
func (o *orderStorage) GetByID(ctx context.Context, id string) (*model.Order, error) {
	ctx, span := tracer.StartTrace(ctx, "psql-storage-get-order-dy-id")
	span.SetAttributes(attribute.String("order-id", id))
	defer span.End()

	query := `
		SELECT 	o.order_uid, o.track_number, o.entry, o.locale, o.internal_signature,
       			o.customer_id, o.delivery_service, o.shardkey, o.sm_id, o.oof_shard,
       			o.date_created, d.name, d.phone, d.zip, d.city, d.address, d.region, d.email,
       			t.id, t.request_id, t.currency, t.provider, t.amount, t.payment_dt,
       			t.bank, t.delivery_cost, t.goods_total, t.custom_fee
		FROM orders o
		JOIN delivery d
    		ON o.order_uid=d.order_uid
		JOIN transaction t
    		ON o.order_uid=t.id
		WHERE o.order_uid=$1	
	`

	var (
		order      model.Order
		createDate time.Time
	)

	err := o.pool.QueryRow(ctx, query, id).Scan(
		&order.OrderUid, &order.TrackNumber, &order.Entry, &order.Locale, &order.InternalSignature, &order.CustomerID, &order.DeliveryService,
		&order.ShardKey, &order.SmID, &order.OofShard, &createDate, &order.Delivery.Name, &order.Delivery.Phone, &order.Delivery.Zip,
		&order.Delivery.City, &order.Delivery.Address, &order.Delivery.Region, &order.Delivery.Email, &order.Payment.Transaction,
		&order.Payment.RequestID, &order.Payment.Currency, &order.Payment.Provider, &order.Payment.Amount, &order.Payment.PaymentDt,
		&order.Payment.Bank, &order.Payment.DeliveryCost, &order.Payment.GoodsTotal, &order.Payment.CustomFee,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return &model.Order{Items: []*model.Product{}}, common.WrapError{Err: domain.ErrOrderNotExists, Msg: domain.ErrOrderNotExists.Error()}
		}

		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case domain.CodeInvalidSyntax:
				return &model.Order{}, common.WrapError{Err: domain.ErrInvalidSyntax, Msg: pgErr.Message}
			}
		}
		return &model.Order{Items: []*model.Product{}}, common.WrapError{Err: err, Msg: "fail to get order by id"}
	}

	order.DateCreated = createDate.String()

	products, err := o.getProductByOrderTrackNumber(ctx, order.TrackNumber)
	if err != nil {
		return &model.Order{Items: []*model.Product{}}, err
	}
	order.Items = products

	return &order, nil
}

// getProductByOrderTrackNumber вернуть items по order track number
func (o *orderStorage) getProductByOrderTrackNumber(ctx context.Context, trackNumber string) ([]*model.Product, error) {
	ctx, span := tracer.StartTrace(ctx, "psql-storage-get-product-by-order-tracknumber")
	span.SetAttributes(attribute.String("tracknumber", trackNumber))
	defer span.End()

	query := `
		SELECT chrt_id, track_number, price, rid, name, sale,
		    size, total_price, nm_id, brand, status
		FROM product
		WHERE track_number=$1
	`

	rows, err := o.pool.Query(ctx, query, trackNumber)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return []*model.Product{}, common.WrapError{Err: domain.ErrItemsNotExists, Msg: domain.ErrItemsNotExists.Error()}
		}
		return []*model.Product{}, common.WrapError{Err: err, Msg: "fail to get items by track number"}
	}

	var products []*model.Product
	for rows.Next() {
		var product model.Product
		err := rows.Scan(&product.ChrtID, &product.TrackNumber, &product.Price,
			&product.Rid, &product.Name, &product.Sale, &product.Size, &product.TotalPrice,
			&product.NmID, &product.Brand, &product.Status,
		)
		if err != nil {
			return []*model.Product{}, common.WrapError{Err: err, Msg: "fail to scan rows"}
		}

		products = append(products, &product)
	}

	return products, nil
}
