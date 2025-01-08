package routes

import (
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5"
)

func RegisterRoutes(router *mux.Router, dbpool *pgx.Conn) {
	api := router.PathPrefix("/api/v1").Subrouter()

	registerBusRoutes(api, dbpool)
}
