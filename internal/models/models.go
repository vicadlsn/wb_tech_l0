package models

import (
	"errors"
	"time"
)

var (
	ErrOrderNotFound = errors.New("order not found")
)

type Order struct {
	OrderUID          string    `json:"order_uid" validate:"required"`
	TrackNumber       string    `json:"track_number" validate:"required"`
	Entry             string    `json:"entry" validate:"required"`
	Locale            string    `json:"locale" validate:"required"`
	InternalSignature *string   `json:"internal_signature"`
	CustomerID        string    `json:"customer_id" validate:"required"`
	DeliveryService   string    `json:"delivery_service" validate:"required"`
	ShardKey          string    `json:"shardkey" validate:"required"`
	SmID              int       `json:"sm_id" validate:"gte=0"`
	DateCreated       time.Time `json:"date_created" validate:"required"`
	OofShard          string    `json:"oof_shard" validate:"required"`

	Delivery Delivery `json:"delivery" validate:"required"`
	Payment  Payment  `json:"payment" validate:"required"`
	Items    []Item   `json:"items" validate:"required"`
}

type Delivery struct {
	OrderUID string `json:"-"`
	Name     string `json:"name" validate:"required"`
	Phone    string `json:"phone" validate:"required,e164"`
	Zip      string `json:"zip" validate:"required"`
	City     string `json:"city" validate:"required"`
	Address  string `json:"address" validate:"required"`
	Region   string `json:"region" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
}

type Payment struct {
	OrderUID     string  `json:"-"`
	Transaction  string  `json:"transaction" validate:"required"`
	RequestID    *string `json:"request_id"`
	Currency     string  `json:"currency" validate:"required"`
	Provider     string  `json:"provider" validate:"required"`
	Amount       int     `json:"amount" validate:"gte=0"`
	PaymentDt    int64   `json:"payment_dt" validate:"gte=0"`
	Bank         string  `json:"bank" validate:"required"`
	DeliveryCost int     `json:"delivery_cost" validate:"gte=0"`
	GoodsTotal   int     `json:"goods_total" validate:"gte=0"`
	CustomFee    int     `json:"custom_fee" validate:"gte=0"`
}

type Item struct {
	ID          int    `json:"-"`
	OrderUID    string `json:"-"`
	ChrtID      int64  `json:"chrt_id" validate:"gte=0"`
	TrackNumber string `json:"track_number" validate:"required"`
	Price       int    `json:"price" validate:"gte=0"`
	Rid         string `json:"rid" validate:"required"`
	Name        string `json:"name" validate:"required"`
	Sale        int    `json:"sale" validate:"gte=0"`
	Size        string `json:"size" validate:"required"`
	TotalPrice  int    `json:"total_price" validate:"gte=0"`
	NmID        int64  `json:"nm_id" validate:"gte=0"`
	Brand       string `json:"brand" validate:"required"`
	Status      int    `json:"status" validate:"gte=0"`
}
