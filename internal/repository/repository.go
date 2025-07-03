package repository

import (
	"context"
	"fmt"

	"webtechl0/internal/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type OrderRepository struct {
	db *pgxpool.Pool
}

func NewOrderRepository(db *pgxpool.Pool) *OrderRepository {
	return &OrderRepository{db: db}
}

func (r *OrderRepository) CreateOrder(ctx context.Context, order *models.Order) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	if err := r.createOrder(ctx, tx, order); err != nil {
		return fmt.Errorf("failed to create order: %w", err)
	}

	if err := r.createDelivery(ctx, tx, &order.Delivery, order.OrderUID); err != nil {
		return fmt.Errorf("failed to create delivery: %w", err)
	}

	if err := r.createPayment(ctx, tx, &order.Payment, order.OrderUID); err != nil {
		return fmt.Errorf("failed to create payment: %w", err)
	}

	if err := r.createItems(ctx, tx, order.Items, order.OrderUID); err != nil {
		return fmt.Errorf("failed to create items: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (r *OrderRepository) createOrder(ctx context.Context, tx pgx.Tx, order *models.Order) error {
	query := `INSERT INTO orders (order_uid, track_number, entry, locale, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard)
              VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`
	_, err := tx.Exec(ctx, query, order.OrderUID, order.TrackNumber, order.Entry, order.Locale, order.CustomerID,
		order.DeliveryService, order.ShardKey, order.SmID, order.DateCreated, order.OofShard)
	return err
}

func (r *OrderRepository) createDelivery(ctx context.Context, tx pgx.Tx, delivery *models.Delivery, orderUID string) error {
	query := `INSERT INTO delivery (order_uid, name, phone, zip, city, address, region, email)
              VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err := tx.Exec(ctx, query, orderUID, delivery.Name, delivery.Phone, delivery.Zip, delivery.City,
		delivery.Address, delivery.Region, delivery.Email)
	return err
}

func (r *OrderRepository) createPayment(ctx context.Context, tx pgx.Tx, payment *models.Payment, orderUID string) error {
	query := `INSERT INTO payment (order_uid, transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee)
              VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`
	_, err := tx.Exec(ctx, query, orderUID, payment.Transaction, payment.RequestID, payment.Currency, payment.Provider,
		payment.Amount, payment.PaymentDt, payment.Bank, payment.DeliveryCost, payment.GoodsTotal, payment.CustomFee)
	return err
}

func (r *OrderRepository) createItems(ctx context.Context, tx pgx.Tx, items []models.Item, orderUID string) error {
	for _, item := range items {
		query := `INSERT INTO item (chrt_id, order_uid, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status)
                  VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`
		_, err := tx.Exec(ctx, query, item.ChrtID, orderUID, item.TrackNumber, item.Price, item.Rid, item.Name,
			item.Sale, item.Size, item.TotalPrice, item.NmID, item.Brand, item.Status)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *OrderRepository) GetOrder(ctx context.Context, orderUID string) (*models.Order, error) {
	var order models.Order
	query := `SELECT order_uid, track_number, entry, locale, internal_signature,customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard FROM orders WHERE order_uid = $1`
	err := r.db.QueryRow(ctx, query, orderUID).Scan(&order.OrderUID, &order.TrackNumber, &order.Entry, &order.Locale, &order.InternalSignature, &order.CustomerID, &order.DeliveryService, &order.ShardKey, &order.SmID, &order.DateCreated, &order.OofShard)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, models.ErrOrderNotFound
		}
		return nil, fmt.Errorf("failed to select order by order_uid: %w", err)
	}

	query = `SELECT name, phone, zip, city, address, region, email FROM delivery WHERE order_uid = $1`
	err = r.db.QueryRow(ctx, query, orderUID).Scan(&order.Delivery.Name, &order.Delivery.Phone, &order.Delivery.Zip, &order.Delivery.City,
		&order.Delivery.Address, &order.Delivery.Region, &order.Delivery.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to select delivery by order_uid: %w", err)
	}

	query = `SELECT transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee 
             FROM payment WHERE order_uid = $1`
	err = r.db.QueryRow(ctx, query, orderUID).Scan(&order.Payment.Transaction, &order.Payment.RequestID, &order.Payment.Currency, &order.Payment.Provider,
		&order.Payment.Amount, &order.Payment.PaymentDt, &order.Payment.Bank, &order.Payment.DeliveryCost, &order.Payment.GoodsTotal, &order.Payment.CustomFee)
	if err != nil {
		return nil, fmt.Errorf("failed to select payment by order_uid: %w", err)
	}

	query = `SELECT chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status FROM item WHERE order_uid = $1`
	rows, err := r.db.Query(ctx, query, orderUID)
	if err != nil {
		return nil, fmt.Errorf("failed to select item by order_uid: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var item models.Item
		err := rows.Scan(&item.ChrtID, &item.TrackNumber, &item.Price, &item.Rid, &item.Name, &item.Sale,
			&item.Size, &item.TotalPrice, &item.NmID, &item.Brand, &item.Status)
		if err != nil {
			return nil, fmt.Errorf("failed to scan item: %w", err)
		}
		order.Items = append(order.Items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return &order, nil
}
