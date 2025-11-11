package product

import (
	"database/sql"
	"github.com/andresdal/gobackexample/types"
)

type Store struct {
	DB *sql.DB
}
func NewStore(db *sql.DB) *Store {
	return &Store{DB: db}
}

func (s *Store) GetProducts() ([]types.Product, error) {
	rows, err := s.DB.Query("SELECT * FROM products")
	if err != nil {
		return nil, err
	}

	products := make([]types.Product, 0)
	for rows.Next() {
		p, err := scanProduct(rows)
		if err != nil {
			return nil, err
		}

		products = append(products, *p)
	}

	return products, nil
}

func scanProduct(row *sql.Rows) (*types.Product, error) {
	p := new(types.Product)
	err := row.Scan(
		&p.ID,
		&p.Name,
		&p.Description,
		&p.Image,
		&p.Price,
		&p.Quantity,
		&p.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return p, nil
}

func (s *Store) CreateProduct(p types.Product) error {
	_, err := s.DB.Exec(
		"INSERT INTO products (name, description, image, price, quantity) VALUES (?, ?, ?, ?, ?)",
		p.Name,
		p.Description,
		p.Image,
		p.Price,
		p.Quantity,
	)
	if err != nil {
		return err
	}
	return nil
}