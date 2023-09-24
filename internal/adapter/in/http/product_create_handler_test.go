package http_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	web "example.com/apiserver/internal/adapter/in/http"
	"example.com/apiserver/internal/application/domain"
	"example.com/apiserver/internal/port/in"
	"github.com/stretchr/testify/assert"
)

func TestCreateProductHandler_Success_WhenCreatingProduct(t *testing.T) {
	uc := &in.MockCreateProductUseCase{}
	uc.MockCreateProduct(func(name string, price float64) (string, error) {
		return "1", nil
	})
	createProductHandler := web.NewCreateProductHandler(uc)

	requestBody := `{
		"Name":  "Chair",
		"Price": 200
	}`

	req, err := http.NewRequest(http.MethodPost, "/product", strings.NewReader(requestBody))
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(createProductHandler.HandleCreateProduct)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)
	assert.JSONEq(t, `{"message": "Product created successfully", "data": {"id": "1"}}`, rr.Body.String())
}

func TestCreateProductHandler_ReturnsError_WhenCreatingProduct(t *testing.T) {
	uc := &in.MockCreateProductUseCase{}
	uc.MockCreateProduct(func(name string, price float64) (string, error) {
		return "", domain.ErrInvalidProductName
	})
	createProductHandler := web.NewCreateProductHandler(uc)

	requestBody := `{
		"Name":  "Chair",
		"Price": 200
	}`

	req, err := http.NewRequest(http.MethodPost, "/product", strings.NewReader(requestBody))
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(createProductHandler.HandleCreateProduct)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.JSONEq(t, `{
		"error": "Invalid Name",
		"details": "Product name cannot be an empty string."
	}`, rr.Body.String())
}

func TestCreateProductHandler_ReturnsError_WithInvalidJsonBody(t *testing.T) {
	uc := &in.MockCreateProductUseCase{}
	uc.MockCreateProduct(func(name string, price float64) (string, error) {
		return "1", nil
	})
	createProductHandler := web.NewCreateProductHandler(uc)

	requestBody := `Not a json body`

	req, err := http.NewRequest(http.MethodPost, "/product", strings.NewReader(requestBody))
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(createProductHandler.HandleCreateProduct)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.JSONEq(t, `{
		"error": "Request body is not a proper JSON", 
		"details": "Your create product request has been received at our end, but we could not understand it due to improper format. Please modify your request following our API spec and try again. If the issue keeps happening, Please open an issue at https://github.com/SyedMHashim/inventory-go/issues"
	}`, rr.Body.String())
}
