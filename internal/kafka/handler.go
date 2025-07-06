package kafka

import (
	"context"
	"encoding/json"
	"log/slog"
	"webtechl0/internal/models"

	"github.com/go-playground/validator/v10"
)

type OrderService interface {
	CreateOrder(ctx context.Context, order *models.Order) error
}

type OrderHandler struct {
	orderService OrderService
	validator    *validator.Validate
	lg           *slog.Logger
}

func NewOrderHandler(orderService OrderService, lg *slog.Logger) *OrderHandler {
	return &OrderHandler{orderService: orderService, lg: lg, validator: validator.New()}
}

func (h *OrderHandler) HandleMessage(ctx context.Context, data []byte) error {
	op := "OrderHandler.HandleMessage"
	var order models.Order
	if err := json.Unmarshal(data, &order); err != nil {
		h.lg.Warn("Invalid JSON", slog.String("op", op), slog.Any("error", err))
		return nil
	}

	lg := h.lg.With("op", op, "order_uid", order.OrderUID)

	if err := h.validator.Struct(order); err != nil {
		lg.Warn("Validation failed", slog.Any("error", err))
		return nil
	}

	if err := h.orderService.CreateOrder(ctx, &order); err != nil {
		lg.Error("Failed to save order", slog.Any("error", err))
		return nil
	}

	lg.Info("Created order")
	return nil
}
