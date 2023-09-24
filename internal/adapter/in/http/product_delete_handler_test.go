package http_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	web "example.com/apiserver/internal/adapter/in/http"
	"example.com/apiserver/internal/port/in"
	"github.com/stretchr/testify/assert"
)

func TestDeleteProductHandler_Success_WhenDeletingProduct(t *testing.T) {
	uc := &in.MockDeleteProductUseCase{}
	uc.MockDeleteProduct(func(id string) error {
		return nil
	})
	deleteProductHandler := web.NewDeleteProductHandler(uc)

	req, err := http.NewRequest(http.MethodDelete, "/product/1", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(deleteProductHandler.HandleDeleteProduct)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.JSONEq(t, `{"message": "Product deleted successfully", "data": {"id": "1"}}`, rr.Body.String())
}

func TestDeleteProductHandler_ReturnsError_WhenDeletingProduct(t *testing.T) {
	uc := &in.MockDeleteProductUseCase{}
	uc.MockDeleteProduct(func(id string) error {
		return errors.New("error deleting product")
	})
	deleteProductHandler := web.NewDeleteProductHandler(uc)

	req, err := http.NewRequest(http.MethodDelete, "/product/1", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(deleteProductHandler.HandleDeleteProduct)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	assert.JSONEq(t, `{
			"error":   "Something wrong happened at our end",
			"details": "Please try again after a while. If the issue keeps happening, Please open an issue at https://github.com/SyedMHashim/inventory-go/issues"
		}`, rr.Body.String())
}
