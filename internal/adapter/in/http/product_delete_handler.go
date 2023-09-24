package http

import (
	"context"
	"net/http"
	"regexp"

	"example.com/apiserver/internal/port/in"
)

type DeleteProductHandler struct {
	uc in.DeleteProductUseCase
}

func NewDeleteProductHandler(uc in.DeleteProductUseCase) *DeleteProductHandler {
	return &DeleteProductHandler{uc: uc}
}

func (h *DeleteProductHandler) HandleDeleteProduct(w http.ResponseWriter, r *http.Request) {
	var id string
	regexCompiler := regexp.MustCompile(`/product/(.*)`)
	match := regexCompiler.FindStringSubmatch(r.URL.Path)
	if len(match) > 1 {
		id = match[1]
	}
	err := h.uc.DeleteProduct(context.Background(), id)
	if err != nil {
		sendError(w, statusCodeOf(err), errorResponseOf(err))
		return
	}
	resp := newSuccessfulResponse(
		"Product deleted successfully",
		responseData{"id": id},
	)
	sendResponse(w, http.StatusOK, resp)
}
