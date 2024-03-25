package http

import (
	"github.com/gorilla/mux"
	"go-micro-services/internal/model/response"
	"go-micro-services/internal/use_case"
	"net/http"
)

type ProductController struct {
	ProductUseCase *use_case.ProductUseCase
}

func NewProductController(productUseCase *use_case.ProductUseCase) *ProductController {
	productControler := &ProductController{
		ProductUseCase: productUseCase,
	}
	return productControler
}
func (ProductController *ProductController) GetOneById(writer http.ResponseWriter, reader *http.Request) {
	vars := mux.Vars(reader)
	id := vars["id"]
	foundProduct, foundProductErr := ProductController.ProductUseCase.GetOneById(id)
	if foundProductErr == nil {
		response.NewResponse(writer, foundProduct)
	}
}
