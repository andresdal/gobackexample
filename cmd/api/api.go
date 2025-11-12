package api

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/andresdal/gobackexample/service/cart"
	"github.com/andresdal/gobackexample/service/order"
	"github.com/andresdal/gobackexample/service/product"
	"github.com/andresdal/gobackexample/service/user"
	"github.com/gorilla/mux"
)

type APIServer struct {
	addr string
	db   *sql.DB
}

func NewAPIServer(addr string, db *sql.DB) *APIServer {
	return &APIServer{
		addr: addr,
		db:   db,
	}
}

func (s *APIServer) Run() error {
	router := mux.NewRouter()
	subrouter := router.PathPrefix("/api/v1").Subrouter()

	userStore := user.NewStore(s.db)
	userHandler := user.NewHandler(userStore)
	userHandler.RegisterRoutes(subrouter)

	productStore := product.NewStore(s.db)
	productHandler := product.NewHandler(productStore)
	productHandler.RegisterRoutes(subrouter)

	orderStore := order.NewStore(s.db)

	cartHandler := cart.NewHandler(orderStore, productStore, userStore)
	cartHandler.RegisterRoutes(subrouter)

	
	log.Println("API server is listening on ", s.addr)

	return http.ListenAndServe(s.addr, router)
}