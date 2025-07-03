package service

import (
	"context"
	"webtechl0/internal/models"
)

type OrderRepository interface {
	CreateOrder(ctx context.Context, order *models.Order) error
	GetOrder(ctx context.Context, orderUID string) (*models.Order, error)
}

type OrderService struct {
	repo OrderRepository
}

func NewOrderService(repo OrderRepository) *OrderService {
	return &OrderService{repo: repo}
}

func (s *OrderService) GetOrder(ctx context.Context, orderUID string) (*models.Order, error) {
	order, err := s.repo.GetOrder(ctx, orderUID)
	if err != nil {
		return nil, err
	}
	return order, nil
}

func (s *OrderService) CreateOrder(ctx context.Context, order *models.Order) error {
	if err := s.repo.CreateOrder(ctx, order); err != nil {
		return err
	}

	return nil
}
