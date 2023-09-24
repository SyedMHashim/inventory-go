package mysql

import (
	"context"
	"database/sql"

	"example.com/apiserver/internal/application/domain"
	_ "github.com/go-sql-driver/mysql"
)

type ProductPersistenceAdapter struct {
	db *sql.DB
}

func NewProductPersistenceAdapter(db *sql.DB) *ProductPersistenceAdapter {
	return &ProductPersistenceAdapter{db: db}
}

func (a *ProductPersistenceAdapter) InsertProduct(ctx context.Context, product domain.Product) error {
	_, err := a.db.ExecContext(ctx, "INSERT INTO products (id, name, price) VALUES(?, ?, ?)", product.ID, product.Name, product.Price)
	if err != nil {
		return err
	}
	return nil
}

func (a *ProductPersistenceAdapter) UpdateProduct(ctx context.Context, product domain.Product) error {
	_, err := a.db.ExecContext(ctx, "UPDATE products SET name = ?, price = ? WHERE id = ?", product.Name, product.Price, product.ID)
	if err != nil {
		return err
	}
	return nil
}

func (a *ProductPersistenceAdapter) GetProductByID(ctx context.Context, id string) (domain.Product, error) {
	p := struct {
		id    string
		name  string
		price float64
	}{}

	rows := a.db.QueryRowContext(ctx, "SELECT * FROM products WHERE id = ?", id)
	err := rows.Scan(&p.id, &p.name, &p.price)
	if err != nil {
		return domain.Product{}, err
	}

	product, err := domain.NewProduct(p.id, p.name, p.price)
	if err != nil {
		return domain.Product{}, err
	}
	return product, nil
}

func (a *ProductPersistenceAdapter) GetAllProducts(ctx context.Context) ([]domain.Product, error) {
	rows, err := a.db.QueryContext(ctx, "SELECT * FROM products")
	if err != nil {
		return nil, err
	}
	var products []domain.Product
	for rows.Next() {
		p := struct {
			id    string
			name  string
			price float64
		}{}

		err := rows.Scan(&p.id, &p.name, &p.price)
		if err != nil {
			return nil, err
		}
		product, err := domain.NewProduct(p.id, p.name, p.price)
		if err != nil {
			return []domain.Product{}, err
		}
		products = append(products, product)
	}
	return products, nil
}

func (a *ProductPersistenceAdapter) DeleteProduct(ctx context.Context, id string) error {
	_, err := a.db.ExecContext(ctx, "DELETE FROM products WHERE id = ?", id)
	if err != nil {
		return err
	}
	return nil
}
