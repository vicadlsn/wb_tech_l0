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

	delivery, err := r.getDelivery(ctx, orderUID)
	if err != nil {
		return nil, err
	}
	order.Delivery = *delivery

	payment, err := r.getPayment(ctx, orderUID)
	if err != nil {
		return nil, err
	}
	order.Payment = *payment

	items, err := r.getItems(ctx, orderUID)
	if err != nil {
		return nil, err
	}
	order.Items = items

	return &order, nil
}

func (r *OrderRepository) getDelivery(ctx context.Context, orderUID string) (*models.Delivery, error) {
	var delivery models.Delivery
	query := `SELECT name, phone, zip, city, address, region, email FROM delivery WHERE order_uid = $1`
	err := r.db.QueryRow(ctx, query, orderUID).Scan(&delivery.Name, &delivery.Phone, &delivery.Zip, &delivery.City,
		&delivery.Address, &delivery.Region, &delivery.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to select delivery by order_uid: %w", err)
	}
	return &delivery, nil
}

func (r *OrderRepository) getPayment(ctx context.Context, orderUID string) (*models.Payment, error) {
	var p models.Payment
	query := `SELECT transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee 
              FROM payment WHERE order_uid = $1`
	err := r.db.QueryRow(ctx, query, orderUID).Scan(&p.Transaction, &p.RequestID, &p.Currency, &p.Provider,
		&p.Amount, &p.PaymentDt, &p.Bank, &p.DeliveryCost, &p.GoodsTotal, &p.CustomFee)
	if err != nil {
		return nil, fmt.Errorf("failed to select payment by order_uid: %w", err)
	}
	return &p, nil
}

func (r *OrderRepository) getItems(ctx context.Context, orderUID string) ([]models.Item, error) {
	query := `SELECT chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status FROM item WHERE order_uid = $1`
	rows, err := r.db.Query(ctx, query, orderUID)
	if err != nil {
		return nil, fmt.Errorf("failed to select items by order_uid: %w", err)
	}
	defer rows.Close()

	var items []models.Item
	for rows.Next() {
		var item models.Item
		err := rows.Scan(&item.ChrtID, &item.TrackNumber, &item.Price, &item.Rid, &item.Name, &item.Sale,
			&item.Size, &item.TotalPrice, &item.NmID, &item.Brand, &item.Status)
		if err != nil {
			return nil, fmt.Errorf("failed to scan item: %w", err)
		}
		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return items, nil
}

func (r *OrderRepository) GetOrders(ctx context.Context) ([]*models.Order, error) {
	query := `SELECT order_uid, track_number, entry, locale, internal_signature,customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard FROM orders`

	rows, err := r.db.Query(ctx, query)

	if err != nil {
		return nil, fmt.Errorf("failed to select orders: %w", err)
	}

	defer rows.Close()

	var orders []*models.Order
	for rows.Next() {
		var order models.Order
		err := rows.Scan(&order.OrderUID, &order.TrackNumber, &order.Entry, &order.Locale, &order.InternalSignature, &order.CustomerID, &order.DeliveryService, &order.ShardKey, &order.SmID, &order.DateCreated, &order.OofShard)

		if err != nil {
			return nil, fmt.Errorf("failed to scan order: %w", err)
		}

		delivery, err := r.getDelivery(ctx, order.OrderUID)
		if err != nil {
			return nil, err
		}
		order.Delivery = *delivery

		payment, err := r.getPayment(ctx, order.OrderUID)
		if err != nil {
			return nil, err
		}
		order.Payment = *payment

		items, err := r.getItems(ctx, order.OrderUID)
		if err != nil {
			return nil, err
		}
		order.Items = items

		orders = append(orders, &order)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return orders, nil
}
