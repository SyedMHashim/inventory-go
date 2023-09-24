package mysql_test

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"testing"

	"example.com/apiserver/internal/adapter/out/mysql"
	"example.com/apiserver/internal/application/domain"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func newSQLMock() (*sql.DB, sqlmock.Sqlmock) {
	db, sqlMock, err := sqlmock.New()
	if err != nil {
		log.Fatal("Error Initialising Mock DB connection")
	}
	return db, sqlMock
}

func TestProductPersistenceAdapter_Success_WhenInsertingProduct(t *testing.T) {
	testdb, sqlMock := newSQLMock()
	sqlMock.ExpectExec("INSERT INTO products").WillReturnResult(sqlmock.NewResult(1, 1))

	a := mysql.NewProductPersistenceAdapter(testdb)

	product := domain.Product{ID: "1", Name: "Chair", Price: 200.00}
	err := a.InsertProduct(context.Background(), product)
	assert.NoError(t, err)
	testdb.Close()
}

func TestProductPersistenceAdapter_ReturnsError_WhenInsertingProduct(t *testing.T) {
	testdb, sqlMock := newSQLMock()
	sqlMock.ExpectExec("INSERT INTO products").WillReturnError(errors.New("error inserting product"))

	a := mysql.NewProductPersistenceAdapter(testdb)

	product := domain.Product{ID: "1", Name: "Chair", Price: 200.00}
	err := a.InsertProduct(context.Background(), product)
	assert.Error(t, err)
	testdb.Close()
}

func TestProductPersistenceAdapter_Success_WhenUpdatingProduct(t *testing.T) {
	testdb, sqlMock := newSQLMock()
	sqlMock.ExpectExec("UPDATE products").WillReturnResult(sqlmock.NewResult(1, 1))

	a := mysql.NewProductPersistenceAdapter(testdb)

	product := domain.Product{ID: "1", Name: "Table", Price: 500}
	err := a.UpdateProduct(context.Background(), product)
	assert.NoError(t, err)
	testdb.Close()
}

func TestProductPersistenceAdapter_ReturnsError_WhenUpdatingProduct(t *testing.T) {
	testdb, sqlMock := newSQLMock()
	sqlMock.ExpectExec("UPDATE products").WillReturnError(errors.New("error updating product"))

	a := mysql.NewProductPersistenceAdapter(testdb)

	product := domain.Product{ID: "1", Name: "Table", Price: 500}
	err := a.UpdateProduct(context.Background(), product)
	assert.Error(t, err)
	testdb.Close()
}

func TestProductPersistenceAdapter_Success_WhenGettingProductByID(t *testing.T) {
	testdb, sqlMock := newSQLMock()
	rows := sqlmock.NewRows([]string{"id", "name", "price"}).
		AddRow("1", "Chair", 200)
	sqlMock.ExpectQuery("SELECT \\* FROM products WHERE id = ?").WillReturnRows(rows)

	a := mysql.NewProductPersistenceAdapter(testdb)

	product, err := a.GetProductByID(context.Background(), "1")
	assert.NoError(t, err)
	assert.Equal(t, "1", product.ID)
	assert.Equal(t, "Chair", product.Name)
	assert.Equal(t, 200.00, product.Price)
	testdb.Close()
}

func TestProductPersistenceAdapter_ReturnsError_WhenGettingProductByID(t *testing.T) {
	testdb, sqlMock := newSQLMock()
	sqlMock.ExpectQuery("SELECT \\* FROM products WHERE id = ?").WillReturnError(errors.New("error getting product by id"))

	a := mysql.NewProductPersistenceAdapter(testdb)

	_, err := a.GetProductByID(context.Background(), "1")
	assert.Error(t, err)
	testdb.Close()
}

func TestProductPersistenceAdapter_Success_WhenGettingAllProducts(t *testing.T) {
	testdb, sqlMock := newSQLMock()
	rows := sqlmock.NewRows([]string{"id", "name", "price"}).
		AddRow(1, "Chair", 200).
		AddRow(2, "Table", 500)
	sqlMock.ExpectQuery("SELECT (.+) FROM products").WillReturnRows(rows)

	a := mysql.NewProductPersistenceAdapter(testdb)

	products, err := a.GetAllProducts(context.Background())
	assert.NoError(t, err)
	assert.Len(t, products, 2)
	assert.Equal(t, "1", products[0].ID)
	assert.Equal(t, "Chair", products[0].Name)
	assert.Equal(t, 200.00, products[0].Price)
	assert.Equal(t, "2", products[1].ID)
	assert.Equal(t, "Table", products[1].Name)
	assert.Equal(t, 500.00, products[1].Price)
	testdb.Close()
}

func TestProductPersistenceAdapter_ReturnsError_WhenGettingAllProducts(t *testing.T) {
	testdb, sqlMock := newSQLMock()
	sqlMock.ExpectQuery("SELECT (.+) FROM products").WillReturnError(errors.New("error getting all products"))

	a := mysql.NewProductPersistenceAdapter(testdb)

	products, err := a.GetAllProducts(context.Background())
	assert.Error(t, err)
	assert.Len(t, products, 0)
	testdb.Close()
}

func TestProductPersistenceAdapter_Success_WhenDeletingProduct(t *testing.T) {
	testdb, sqlMock := newSQLMock()
	sqlMock.ExpectExec("DELETE FROM products").WillReturnResult(sqlmock.NewResult(1, 1))

	a := mysql.NewProductPersistenceAdapter(testdb)

	err := a.DeleteProduct(context.Background(), "1")
	assert.NoError(t, err)
	testdb.Close()
}

func TestProductPersistenceAdapter_ReturnsError_WhenDeletingProduct(t *testing.T) {
	testdb, sqlMock := newSQLMock()
	sqlMock.ExpectExec("DELETE FROM products").WillReturnError(errors.New("error deleting product"))

	a := mysql.NewProductPersistenceAdapter(testdb)

	err := a.DeleteProduct(context.Background(), "1")
	assert.Error(t, err)
	testdb.Close()
}
