package service

import (
	"context"
	"fmt"
	"log/slog"
	"webtechl0/internal/models"
)

type OrderRepository interface {
	CreateOrder(ctx context.Context, order *models.Order) error
	GetOrder(ctx context.Context, orderUID string) (*models.Order, error)
	GetOrders(ctx context.Context) ([]*models.Order, error)
}

type OrderCache interface {
	Get(orderUID string) (*models.Order, bool)
	Put(orderUID string, order *models.Order)
}

type OrderService struct {
	repo  OrderRepository
	cache OrderCache
	lg    *slog.Logger
}

func NewOrderService(repo OrderRepository, cache OrderCache, lg *slog.Logger) *OrderService {
	return &OrderService{repo: repo, cache: cache, lg: lg}
}

func (s *OrderService) GetOrder(ctx context.Context, orderUID string) (*models.Order, error) {
	if order, ok := s.cache.Get(orderUID); ok {
		s.lg.Debug("Got order from cache", slog.Any("order_uid", orderUID))
		return order, nil
	}

	order, err := s.repo.GetOrder(ctx, orderUID)
	if err != nil {
		return nil, err
	}

	s.lg.Debug("Got order from DB", slog.Any("order_uid", orderUID))

	s.cache.Put(orderUID, order)
	return order, nil
}

func (s *OrderService) GetOrders(ctx context.Context) ([]*models.Order, error) {
	orders, err := s.repo.GetOrders(ctx)
	if err != nil {
		return nil, err
	}

	return orders, nil
}

func (s *OrderService) CreateOrder(ctx context.Context, order *models.Order) error {
	if err := s.repo.CreateOrder(ctx, order); err != nil {
		return err
	}

	return nil
}

func (s *OrderService) FillCache(ctx context.Context) error {
	orders, err := s.repo.GetOrders(ctx)

	if err != nil {
		return fmt.Errorf("failed to load orders from DB: %w", err)
	}

	for _, order := range orders {
		s.cache.Put(order.OrderUID, order)

	}

	return nil
}
