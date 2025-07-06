package handler

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"webtechl0/internal/models"
)

type OrderService interface {
	GetOrder(ctx context.Context, orderUID string) (*models.Order, error)
	GetOrders(ctx context.Context) ([]*models.Order, error)
}

type OrderHandler struct {
	orderService OrderService
	lg           *slog.Logger
}

func NewOrderHandler(orderService OrderService, lg *slog.Logger) *OrderHandler {
	return &OrderHandler{orderService: orderService, lg: lg}
}

func (h *OrderHandler) GetOrder(w http.ResponseWriter, r *http.Request) {
	op := "OrderHandler.GetOrder"
	orderUID := r.PathValue("order_uid")
	log := h.lg.With(slog.String("op", op), slog.String("order_uid", orderUID))

	order, err := h.orderService.GetOrder(r.Context(), orderUID)

	if err != nil {
		if errors.Is(err, models.ErrOrderNotFound) {
			log.Info("Order not found")
			http.Error(w, "Order not found", http.StatusNotFound)
			return
		}

		log.Error("Internal server error", slog.Any("error", err))
		http.Error(w, "Internal server error", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(order); err != nil {
		log.Error("Failed to encode response", slog.Any("error", err))
	}
}

func (h *OrderHandler) GetAllOrders(w http.ResponseWriter, r *http.Request) {
	op := "OrderHandler.GetAllOrders"
	log := h.lg.With(slog.String("op", op))

	orders, err := h.orderService.GetOrders(r.Context())

	if err != nil {
		log.Error("Internal server error", slog.Any("error", err))
		http.Error(w, "Internal server error", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(orders); err != nil {
		log.Error("Failed to encode response", slog.Any("error", err))
	}
}
