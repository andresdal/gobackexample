package cart

import (
	"fmt"

	"github.com/andresdal/gobackexample/types"
)

func getCartItemsIDs(items []types.CartItem) ([]int, error) {
	productIDs := make([]int, 0, len(items))
	for _, item := range items {
		if item.Quantity <= 0 {
			return nil, fmt.Errorf("invalid quantity for product ID %d", item.ProductID)
		}
		productIDs = append(productIDs, item.ProductID)
	}
	return productIDs, nil
}

func (h *Handler) createOrder(products []types.Product, items []types.CartItem, userID int) (int, float64, error) {
	productMap := make(map[int]types.Product) // para mayor eficiencia
	for _, p := range products {
		productMap[p.ID] = p
	}

	// check if all products are in stock
	if err := checkCartStock(items, productMap); err != nil {
		return 0, 0, err
	}

	// calculate total price
	totalPrice := calculateTotalPrice(items, productMap)

	// reduce stock in DB
	for _, item := range items {
		product := productMap[item.ProductID]
		newQuantity := product.Quantity - item.Quantity
		err := h.productStore.UpdateProductQuantity(item.ProductID, newQuantity)
		if err != nil {
			return 0, 0, err
		}
	}

	// create order
	orderID, _ := h.store.CreateOrder(types.Order{
		UserID: userID,
		Total:  totalPrice,
		Status: "pending",
		Address: "123 Main St", // this would come from the user in a real app
	})

	// create order items
	for _, item := range items {
		h.store.CreateOrderItem(types.OrderItem{
			OrderID:   orderID,
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     productMap[item.ProductID].Price,
		})
	}

	return orderID, totalPrice, nil
}

func checkCartStock(items []types.CartItem, productMap map[int]types.Product) error {
	if len(items) == 0 {
		return fmt.Errorf("cart is empty")
	}

	for _, item := range items {
		product, exists := productMap[item.ProductID]
		if !exists {
			return fmt.Errorf("product ID %d not found", item.ProductID)
		}
		if product.Quantity < item.Quantity {
			return fmt.Errorf("insufficient stock for product %v", product.Name)
		}
	}
	return nil
}

func calculateTotalPrice(items []types.CartItem, productMap map[int]types.Product) float64 {
	total := 0.0
	for _, item := range items {
		product := productMap[item.ProductID]
		total += float64(item.Quantity) * product.Price
	}
	return total
}