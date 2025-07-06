package main

import (
	"context"
	"encoding/json"
	"log/slog"
	"math/rand"
	"os"
	"time"

	"webtechl0/internal/config"

	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
)

const (
	messageCount = 50
)

func main() {
	log := slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	)

	cfg, err := config.New("")
	if err != nil {
		log.Error("Failed to load config", slog.Any("error", err))
		os.Exit(1)
	}

	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers:  cfg.Kafka.Brokers,
		Topic:    cfg.Kafka.Topic,
		Balancer: &kafka.RoundRobin{},
	})
	defer writer.Close()

	ctx := context.Background()
	var msg kafka.Message

	for i := range messageCount {
		orderUID := uuid.NewString()
		if i%10 == 0 {
			log.Info("sending repeated message")
		} else if i%9 == 0 {
			msg = createInvalidMsg()
			log.Info("sending invalid message")
		} else {
			msg = createOrderMsg(orderUID)
			log.Info("sending valid message")
		}

		if err := writer.WriteMessages(ctx, msg); err != nil {
			log.Error("Write error", slog.Any("error", err))
		}

		time.Sleep(300 * time.Millisecond)
	}
}

func createOrderMap(orderUID string) map[string]any {
	itemsCount := rand.Intn(10) + 1
	items := make([]map[string]any, 0, itemsCount)
	item := map[string]any{
		"chrt_id":      9934930,
		"track_number": "WBILMTESTTRACK",
		"price":        453,
		"rid":          "ab4219087a764ae0btest",
		"name":         "Mascaras",
		"sale":         30,
		"size":         "0",
		"total_price":  317,
		"nm_id":        2389212,
		"brand":        "Vivienne Sabo",
		"status":       202,
	}
	for range itemsCount {
		items = append(items, item)
	}

	return map[string]any{
		"order_uid":    orderUID,
		"track_number": "WBILMTESTTRACK",
		"entry":        "WBIL",
		"delivery": map[string]string{
			"name":    "Test Testov",
			"phone":   "+9720000000",
			"zip":     "2639809",
			"city":    "Kiryat Mozkin",
			"address": "Ploshad Mira 15",
			"region":  "Kraiot",
			"email":   "test@gmail.com",
		},
		"payment": map[string]any{
			"transaction":   orderUID,
			"request_id":    "",
			"currency":      "USD",
			"provider":      "wbpay",
			"amount":        1817,
			"payment_dt":    1637907727,
			"bank":          "alpha",
			"delivery_cost": 1500,
			"goods_total":   317,
			"custom_fee":    0,
		},
		"items":              items,
		"locale":             "en",
		"internal_signature": "",
		"customer_id":        "test",
		"delivery_service":   "meest",
		"shardkey":           "9",
		"sm_id":              99,
		"date_created":       "2021-11-26T06:22:19Z",
		"oof_shard":          "1",
	}
}

func createOrderMapInvalidFields(orderUID string) map[string]any {
	return map[string]any{
		"order_uid":    orderUID,
		"track_number": "WBILMTESTTRACK",
		"entry":        "WBIL",
		"delivery": map[string]string{
			"name":    "Test Testov",
			"phone":   "720000000",
			"zip":     "2639809",
			"city":    "Kiryat Mozkin",
			"address": "Ploshad Mira 15",
			"region":  "Kraiot",
			"email":   "invalid@gma",
		},
		"payment": map[string]any{
			"transaction":   orderUID,
			"request_id":    "",
			"currency":      "USD",
			"provider":      "wbpay",
			"amount":        1817,
			"payment_dt":    1637907727,
			"bank":          "alpha",
			"delivery_cost": 1500,
			"custom_fee":    0,
		},
		"items": []map[string]any{
			{
				"chrt_id":      9934930,
				"track_number": "WBILMTESTTRACK",
				"price":        453,
				"rid":          "ab4219087a764ae0btest",
				"name":         "Mascaras",
				"sale":         30,
				"size":         "0",
				"total_price":  317,
				"nm_id":        2389212,
				"brand":        "Vivienne Sabo",
				"status":       202,
			},
		},
		"locale":             "en",
		"internal_signature": "",
		"customer_id":        "test",
		"shardkey":           "9",
		"sm_id":              99,
		"date_created":       "2021-11-26T06:22:19Z",
		"oof_shard":          "1",
	}
}

func createOrderMsg(orderUID string) kafka.Message {
	data, _ := json.Marshal(createOrderMap(orderUID))

	return kafka.Message{
		Key:   []byte(orderUID),
		Value: data,
	}
}

func createInvalidMsg() kafka.Message {
	var (
		data []byte
		key  string
	)

	p := rand.Float64()
	if p < 0.1 {
		key = "empty body"
		data, _ = json.Marshal("")
	} else if p < 0.4 {
		key = "bad json"
		data = []byte(`{"order_uid":`)
	} else {
		key = "invalid fields"
		data, _ = json.Marshal(createOrderMapInvalidFields(key))
	}

	return kafka.Message{
		Key:   []byte(key),
		Value: data,
	}
}
