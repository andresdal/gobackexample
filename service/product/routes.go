package product

import (
	"net/http"

	"github.com/andresdal/gobackexample/types"
	"github.com/andresdal/gobackexample/utils"
	"github.com/gorilla/mux"
)

type Handler struct {
	store types.ProductStore
}

func NewHandler(store types.ProductStore) *Handler {
	return &Handler{store: store}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/products", h.handleGetProducts).Methods("GET")
	router.HandleFunc("/products", h.handleCreateProduct).Methods("POST")
}

func (h *Handler) handleGetProducts(w http.ResponseWriter, r *http.Request) {
	products, err := h.store.GetProducts()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	
	utils.WriteJSON(w, http.StatusOK, products)
}

func (h *Handler) handleCreateProduct(w http.ResponseWriter, r *http.Request) {
	// parse json body
	var p types.Product
	if err := utils.ParseJSON(r, &p); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	// create product
	if err := h.store.CreateProduct(p); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJSON(w, http.StatusCreated, p)
}
