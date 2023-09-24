package http

import (
	"context"
	"encoding/json"
	"net/http"

	"example.com/apiserver/internal/port/in"
)

type CreateProductHandler struct {
	uc in.CreateProductUseCase
}

func NewCreateProductHandler(uc in.CreateProductUseCase) *CreateProductHandler {
	return &CreateProductHandler{uc: uc}
}

func (h *CreateProductHandler) HandleCreateProduct(w http.ResponseWriter, r *http.Request) {
	req := struct {
		Name  string
		Price float64
	}{}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		resp := newErrorResponse(
			"Request body is not a proper JSON",
			"Your create product request has been received at our end, but we could not understand it due to improper format. Please modify your request following our API spec and try again. If the issue keeps happening, Please open an issue at https://github.com/SyedMHashim/inventory-go/issues",
		)
		sendError(w, http.StatusBadRequest, resp)
		return
	}
	id, err := h.uc.CreateProduct(context.Background(), req.Name, req.Price)
	if err != nil {
		sendError(w, statusCodeOf(err), errorResponseOf(err))
		return
	}
	resp := newSuccessfulResponse(
		"Product created successfully",
		responseData{"id": id},
	)
	sendResponse(w, http.StatusCreated, resp)
}
