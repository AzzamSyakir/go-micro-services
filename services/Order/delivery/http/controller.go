package http

import (
	"encoding/json"
	"github.com/gorilla/mux"
	model_request "go-micro-services/services/Product/model/request/controller"
	"go-micro-services/services/Product/model/response"
	"go-micro-services/services/Product/use_case"
	"net/http"
)

type ProductController struct {
	ProductUseCase *use_case.ProductUseCase
}

func NewProductController(productUseCase *use_case.ProductUseCase) *ProductController {
	productController := &ProductController{
		ProductUseCase: productUseCase,
	}
	return productController
}
func (ProductController *ProductController) GetOneById(writer http.ResponseWriter, reader *http.Request) {
	vars := mux.Vars(reader)
	id := vars["id"]
	foundProduct, foundProductErr := ProductController.ProductUseCase.GetOneById(id)
	if foundProductErr == nil {
		response.NewResponse(writer, foundProduct)
	}
}
func (ProductController *ProductController) PatchOneById(writer http.ResponseWriter, reader *http.Request) {
	vars := mux.Vars(reader)
	id := vars["id"]

	request := &model_request.ProductPatchOneByIdRequest{}
	decodeErr := json.NewDecoder(reader.Body).Decode(request)
	if decodeErr != nil {
		panic(decodeErr)
	}
	result := ProductController.ProductUseCase.PatchOneByIdFromRequest(id, request)

	response.NewResponse(writer, result)
}
