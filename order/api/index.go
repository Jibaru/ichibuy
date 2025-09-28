package handler

import (
	"net/http"

	"ichibuy/order/config"
	"ichibuy/order/db"
	_ "ichibuy/order/docs"
	"ichibuy/order/server"
)

var cfg = config.Load()

// Handler for vercel function
func Handler(w http.ResponseWriter, r *http.Request) {
	db, err := db.New(cfg.PostgresURI)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()
	server.New(cfg, db).ServeHTTP(w, r)
}
