package handler

import "net/http"

func NewRouter(orderHandler *OrderHandler) http.Handler {
	mux := http.NewServeMux()

	mux.Handle("/", http.FileServer(http.Dir("./web")))

	mux.HandleFunc("GET /order/{order_uid}/", orderHandler.GetOrder)

	return mux
}
