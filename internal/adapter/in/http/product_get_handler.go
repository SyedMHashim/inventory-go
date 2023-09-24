package http

import (
	"context"
	"net/http"
	"regexp"

	"example.com/apiserver/internal/port/in"
)

type GetProductHandler struct {
	uc in.GetProductUseCase
}

func NewGetProductHandler(uc in.GetProductUseCase) *GetProductHandler {
	return &GetProductHandler{uc: uc}
}

func (h *GetProductHandler) HandleGetProduct(w http.ResponseWriter, r *http.Request) {
	var id string
	regexCompiler := regexp.MustCompile(`/product/(.*)`)
	match := regexCompiler.FindStringSubmatch(r.URL.Path)
	if len(match) > 1 {
		id = match[1]
	}

	product, err := h.uc.GetProduct(context.Background(), id)
	if err != nil {
		sendError(w, statusCodeOf(err), errorResponseOf(err))
		return
	}
	resp := newSuccessfulResponse(
		"Product retrieved successfully",
		responseData{"product": product},
	)
	sendResponse(w, http.StatusOK, resp)
}

func (h *GetProductHandler) HandleGetProducts(w http.ResponseWriter, r *http.Request) {
	products, err := h.uc.GetProducts(context.Background())
	if err != nil {
		sendError(w, statusCodeOf(err), errorResponseOf(err))
		return
	}
	resp := newSuccessfulResponse(
		"Products retrieved successfully",
		responseData{"products": products},
	)
	sendResponse(w, http.StatusOK, resp)
}
