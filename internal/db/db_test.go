package db

import (
	"database/sql"
	"log"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

var (
	testdb  *sql.DB
	sqlMock sqlmock.Sqlmock
)

func setup() {
	var err error
	testdb, sqlMock, err = sqlmock.New()
	if err != nil {
		log.Fatal("Error Initialising Mock DB connection")
	}
}

func tearDown() {
	testdb.Close()
}

func TestMain(m *testing.M) {
	setup()
	m.Run()
	tearDown()
}

func TestCreateProduct(t *testing.T) {
	sqlMock.ExpectExec("INSERT INTO products").WillReturnResult(sqlmock.NewResult(1, 1))

	product := Product{Name: "Chair", Price: 200.00}
	err := CreateProduct(testdb, &product)
	assert.NoError(t, err)
	assert.Equal(t, 1, product.Id)
	assert.Equal(t, "Chair", product.Name)
	assert.Equal(t, 200.00, product.Price)
}

func TestGetProducts(t *testing.T) {
	rows := sqlmock.NewRows([]string{"id", "name", "price"}).
		AddRow(1, "Chair", 200).
		AddRow(2, "Table", 500)
	sqlMock.ExpectQuery("SELECT (.+) FROM products").WillReturnRows(rows)
	products, err := GetProducts(testdb)
	assert.NoError(t, err)
	assert.Len(t, products, 2)
	assert.Equal(t, 1, products[0].Id)
	assert.Equal(t, "Chair", products[0].Name)
	assert.Equal(t, 200.00, products[0].Price)
	assert.Equal(t, 2, products[1].Id)
	assert.Equal(t, "Table", products[1].Name)
	assert.Equal(t, 500.00, products[1].Price)
}

func TestGetProduct(t *testing.T) {
	rows := sqlmock.NewRows([]string{"id", "name", "price"}).
		AddRow(1, "Chair", 200)
	sqlMock.ExpectQuery("SELECT \\* FROM products WHERE id = 1").WillReturnRows(rows)
	product, err := GetProduct(testdb, 1)
	assert.NoError(t, err)
	assert.Equal(t, 1, product.Id)
	assert.Equal(t, "Chair", product.Name)
	assert.Equal(t, 200.00, product.Price)
}

func TestUpdateProduct(t *testing.T) {
	sqlMock.ExpectExec("UPDATE products").WillReturnResult(sqlmock.NewResult(1, 1))
	product := Product{Name: "Table", Price: 500}
	err := UpdateProduct(testdb, product)
	assert.NoError(t, err)
}

func TestDeleteProduct(t *testing.T) {
	sqlMock.ExpectExec("DELETE FROM products").WillReturnResult(sqlmock.NewResult(1, 1))
	err := DeleteProduct(testdb, 1)
	assert.NoError(t, err)
}
