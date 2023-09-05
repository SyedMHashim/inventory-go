package router

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"example.com/apiserver/internal/db"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

var (
	r       *mux.Router
	testdb  *sql.DB
	sqlMock sqlmock.Sqlmock
)

func setup() {
	var err error
	testdb, sqlMock, err = sqlmock.New()
	if err != nil {
		log.Fatal("Error Initialising Mock DB connection")
	}
	r = Initialise(testdb)
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

	product := db.Product{Name: "Chair", Price: 200.00}
	reqBody, _ := json.Marshal(product)
	req, err := http.NewRequest(http.MethodPost, "/product", bytes.NewBuffer(reqBody))
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(createProduct)
	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusCreated, rr.Code)

	respBody, _ := io.ReadAll(rr.Body)
	productCreated := db.Product{}
	_ = json.Unmarshal(respBody, &productCreated)
	assert.Equal(t, 1, productCreated.Id)
	assert.Equal(t, "Chair", productCreated.Name)
	assert.Equal(t, 200.00, productCreated.Price)
}

func TestGetProducts(t *testing.T) {
	rows := sqlmock.NewRows([]string{"id", "name", "price"}).
		AddRow(1, "Chair", 200).
		AddRow(2, "Table", 500)
	sqlMock.ExpectQuery("SELECT (.+) FROM products").WillReturnRows(rows)

	req, err := http.NewRequest(http.MethodGet, "/products", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(getProducts)
	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)

	respBody, _ := io.ReadAll(rr.Body)
	products := make([]db.Product, 2)
	_ = json.Unmarshal(respBody, &products)
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

	req, err := http.NewRequest(http.MethodGet, "/product/1", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(getProduct)
	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)

	respBody, _ := io.ReadAll(rr.Body)
	var product db.Product
	_ = json.Unmarshal(respBody, &product)
	assert.Equal(t, 1, product.Id)
	assert.Equal(t, "Chair", product.Name)
	assert.Equal(t, 200.00, product.Price)
}

func TestUpdateProduct(t *testing.T) {
	rows := sqlmock.NewRows([]string{"id", "name", "price"}).
		AddRow(1, "Chair", 200)
	sqlMock.ExpectQuery("SELECT \\* FROM products WHERE id = 1").WillReturnRows(rows)
	sqlMock.ExpectExec("UPDATE products").WillReturnResult(sqlmock.NewResult(1, 1))

	reqBody := []byte(`{"Name":"Table","Price":500}`)
	req, err := http.NewRequest(http.MethodPut, "/product/1", bytes.NewBuffer(reqBody))
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(updateProduct)
	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)

	respBody, _ := io.ReadAll(rr.Body)
	var product db.Product
	_ = json.Unmarshal(respBody, &product)
	assert.Equal(t, 1, product.Id)
	assert.Equal(t, "Table", product.Name)
	assert.Equal(t, 500.00, product.Price)
}

func TestDeleteProduct(t *testing.T) {
	sqlMock.ExpectExec("DELETE FROM products").WillReturnResult(sqlmock.NewResult(1, 1))

	req, err := http.NewRequest(http.MethodDelete, "/product/1", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(deleteProduct)
	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)

	respBody, _ := io.ReadAll(rr.Body)
	resp := make(map[string]string)
	_ = json.Unmarshal(respBody, &resp)
	assert.Equal(t, "successfully deleted", resp["message"])
}
