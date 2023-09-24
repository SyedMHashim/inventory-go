package http

import (
	"context"
	"encoding/json"
	"net/http"
	"regexp"

	"example.com/apiserver/internal/port/in"
)

type UpdateProductHandler struct {
	uc in.UpdateProductUseCase
}

func NewUpdateProductHandler(uc in.UpdateProductUseCase) *UpdateProductHandler {
	return &UpdateProductHandler{uc: uc}
}

func (h *UpdateProductHandler) HandleUpdateProduct(w http.ResponseWriter, r *http.Request) {
	var id string
	regexCompiler := regexp.MustCompile(`/product/(.*)`)
	match := regexCompiler.FindStringSubmatch(r.URL.Path)
	if len(match) > 1 {
		id = match[1]
	}

	req := struct {
		Name  string
		Price float64
	}{}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		resp := newErrorResponse(
			"Request body is not a proper JSON",
			"Your update product request has been received at our end, but we could not understand it due to improper format. Please modify your request following our API spec and try again. If the issue keeps happening, Please open an issue at https://github.com/SyedMHashim/inventory-go/issues",
		)
		sendError(w, http.StatusBadRequest, resp)
		return
	}
	err = h.uc.UpdateProduct(context.Background(), id, req.Name, req.Price)
	if err != nil {
		sendError(w, statusCodeOf(err), errorResponseOf(err))
		return
	}
	resp := newSuccessfulResponse(
		"Product updated successfully",
		responseData{"id": id},
	)
	sendResponse(w, http.StatusOK, resp)
}
