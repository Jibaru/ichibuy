package handler

import (
	"net/http"

	"ichibuy/auth/config"
	"ichibuy/auth/db"
	_ "ichibuy/auth/docs"
	"ichibuy/auth/server"
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
