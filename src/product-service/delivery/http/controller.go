package http

import (
	"encoding/json"
	model_request "go-micro-services/src/product-service/model/request/controller"
	"go-micro-services/src/product-service/model/response"
	"go-micro-services/src/product-service/use_case"
	"net/http"

	"github.com/gorilla/mux"
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
func (ProductController *ProductController) GetProduct(writer http.ResponseWriter, reader *http.Request) {
	vars := mux.Vars(reader)
	id := vars["id"]
	foundProduct := ProductController.ProductUseCase.GetOneById(id)
	response.NewResponse(writer, foundProduct)
}
func (ProductController *ProductController) UpdateProduct(writer http.ResponseWriter, reader *http.Request) {
	vars := mux.Vars(reader)
	id := vars["id"]

	request := &model_request.ProductPatchOneByIdRequest{}
	decodeErr := json.NewDecoder(reader.Body).Decode(request)
	if decodeErr != nil {
		panic(decodeErr)
	}
	result := ProductController.ProductUseCase.UpdateProduct(id, request)

	response.NewResponse(writer, result)
}

func (productController *ProductController) CreateProduct(writer http.ResponseWriter, reader *http.Request) {

	request := &model_request.CreateProduct{}

	decodeErr := json.NewDecoder(reader.Body).Decode(request)
	if decodeErr != nil {
		http.Error(writer, "Failed to decode request body: "+decodeErr.Error(), http.StatusBadRequest)
		return
	}

	result, _ := productController.ProductUseCase.CreateProduct(request)

	response.NewResponse(writer, result)
}

func (productController *ProductController) ListProduct(writer http.ResponseWriter, reader *http.Request) {
	product := productController.ProductUseCase.ListProduct()
	response.NewResponse(writer, product)
}

func (productController *ProductController) DeleteProduct(writer http.ResponseWriter, reader *http.Request) {
	vars := mux.Vars(reader)
	id := vars["id"]

	result := productController.ProductUseCase.DeleteProduct(id)

	response.NewResponse(writer, result)
}
