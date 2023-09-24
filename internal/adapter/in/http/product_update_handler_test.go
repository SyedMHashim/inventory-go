package http_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	web "example.com/apiserver/internal/adapter/in/http"
	"example.com/apiserver/internal/port/in"
	"github.com/stretchr/testify/assert"
)

func TestUpdateProductHandler_Success_WhenUpdatingProduct(t *testing.T) {
	uc := &in.MockUpdateProductUseCase{}
	uc.MockUpdateProduct(func(id, name string, price float64) error {
		return nil
	})
	updateProductHandler := web.NewUpdateProductHandler(uc)

	requestBody := `{
		"Name":  "Chair",
		"Price": 200
	}`

	req, err := http.NewRequest(http.MethodPut, "/product/1", strings.NewReader(requestBody))
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(updateProductHandler.HandleUpdateProduct)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.JSONEq(t, `{"message": "Product updated successfully", "data": {"id": "1"}}`, rr.Body.String())
}

func TestUpdateProductHandler_ReturnsError_WhenUpdatingProduct(t *testing.T) {
	uc := &in.MockUpdateProductUseCase{}
	uc.MockUpdateProduct(func(id, name string, price float64) error {
		return errors.New("error updating product")
	})
	updateProductHandler := web.NewUpdateProductHandler(uc)

	requestBody := `{
		"Name":  "Chair",
		"Price": 200
	}`

	req, err := http.NewRequest(http.MethodPut, "/product/1", strings.NewReader(requestBody))
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(updateProductHandler.HandleUpdateProduct)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	assert.JSONEq(t, `{
		"error":   "Something wrong happened at our end",
		"details": "Please try again after a while. If the issue keeps happening, Please open an issue at https://github.com/SyedMHashim/inventory-go/issues"
	}`, rr.Body.String())
}

func TestUpdateProductHandler_ReturnsError_WithInvalidJsonBody(t *testing.T) {
	uc := &in.MockUpdateProductUseCase{}
	uc.MockUpdateProduct(func(id, name string, price float64) error {
		return nil
	})
	updateProductHandler := web.NewUpdateProductHandler(uc)

	requestBody := `Not a json body`

	req, err := http.NewRequest(http.MethodPut, "/product/1", strings.NewReader(requestBody))
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(updateProductHandler.HandleUpdateProduct)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.JSONEq(t, `{
		"error": "Request body is not a proper JSON",
		"details": "Your update product request has been received at our end, but we could not understand it due to improper format. Please modify your request following our API spec and try again. If the issue keeps happening, Please open an issue at https://github.com/SyedMHashim/inventory-go/issues"
	}`, rr.Body.String())
}
