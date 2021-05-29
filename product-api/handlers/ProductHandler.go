package handlers

import (
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/mux"
	protos "github.com/kjunn2000/go-microservice-tutorial/currency-server/protos/currency"
	product "github.com/kjunn2000/go-microservice-tutorial/product-api/data"
)

type ProductHandler struct {
	sync.Mutex
	products product.Products
	protos   protos.CurrencyClient
}

func GetProductHandler() *ProductHandler {
	return &ProductHandler{
		products: product.Products{
			product.Product{1, "bottle", 14, time.Now()},
			product.Product{2, "mic", 25, time.Now()},
			product.Product{3, "coffee", 422, time.Now()},
		},
	}
}

// swagger:route GET /products/{id} products listSingle
// Return single products from the database
// responses:
//	200: productResponse
//	404: errorResponse
func (ph *ProductHandler) GetOne(w http.ResponseWriter, request *http.Request) {

	vars := mux.Vars(request)
	productId, _ := strconv.Atoi(vars["id"])
	product, exist := ph.findProductById(productId)
	if !exist {
		return
	}
	ph.toJSON(w, product)
}

// swagger:route GET /products products listProducts
// Return a list of products from the database
// responses:
//	200: productsResponse
func (ph *ProductHandler) Get(w http.ResponseWriter, request *http.Request) {
	fmt.Println("Hello")
	ph.toJSON(w, ph.products)
}

// swagger:route POST /products products createProduct
// Create a new product
// responses:
//	200: productResponse
//  422: errorValidation
//  501: errorResponse
func (ph *ProductHandler) Post(w http.ResponseWriter, r *http.Request) {
	product := r.Context().Value(Product{}).(product.Product)
	ph.products = append(ph.products, product)
}

// swagger:route DELETE /products/{id} products deleteProduct
// Delete a new product
// responses:
//	201: noContentResponse
//  404: errorResponse
//  501: errorResponse
func (ph *ProductHandler) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	res, _ := strconv.Atoi(id)
	index := ph.findProductIndexById(res)
	if index == -1 {
		return
	}
	ph.products = append(ph.products[:index], ph.products[index+1:]...)
	fmt.Fprintln(w, "Deleted")
}

// swagger:route PUT /products products updateProduct
// Update a new product
// responses:
//	200: productResponse
//  404: errorResponse
//  422: errorValidatione
func (ph *ProductHandler) Put(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	product := r.Context().Value(Product{}).(product.Product)
	convId, _ := strconv.Atoi(id)
	index := ph.findProductIndexById(convId)
	if index == -1 {
		return
	}
	ph.products[index].Name = product.Name
	ph.products[index].CreatedAt = product.CreatedAt
}

func (ph *ProductHandler) findProductIndexById(id int) (index int) {
	for i, v := range ph.products {
		if v.Id == id {
			return i
		}
	}
	return -1
}

func (ph *ProductHandler) findProductById(id int) (prodcut product.Product, exist bool) {
	for _, product := range ph.products {
		if product.Id == id {
			return product, true
		}
	}
	return product.Product{}, false
}

func (ph *ProductHandler) toJSON(w http.ResponseWriter, target interface{}) {
	encode := json.NewEncoder(w)
	encode.Encode(target)
}

type Product struct{}

func (ph *ProductHandler) LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.RequestURI)
		var product product.Product
		err := json.NewDecoder(r.Body).Decode(&product)
		if err != nil {
			http.Error(w, "Cannot decode", http.StatusBadRequest)
			return
		}
		err = product.Validate()
		if err != nil {
			http.Error(w, "Incorrect format : "+err.Error(), http.StatusBadRequest)
			return
		}
		ctx := context.WithValue(r.Context(), Product{}, product)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

// ErrInvalidProductPath is an error message when the product path is not valid
var ErrInvalidProductPath = fmt.Errorf("Invalid Path, path should be /products/[id]")

// GenericError is a generic error message returned by a server
type GenericError struct {
	Message string `json:"message"`
}

// ValidationError is a collection of validation error messages
type ValidationError struct {
	Messages []string `json:"messages"`
}

func (ph *ProductHandler) GzipMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			ww := NewWrapperdWriter(w)
			next.ServeHTTP(ww, r)
		}
		next.ServeHTTP(w, r)

	})
}

type WrappedWriter struct {
	rw http.ResponseWriter
	gz *gzip.Writer
}

func NewWrapperdWriter(w http.ResponseWriter) *WrappedWriter {
	return &WrappedWriter{rw: w, gz: gzip.NewWriter(w)}
}

func (ww *WrappedWriter) Header() http.Header {
	return ww.rw.Header()
}

func (ww *WrappedWriter) Flush() {
	ww.gz.Flush()
	ww.gz.Close()
}

func (ww *WrappedWriter) Write(data []byte) (int, error) {
	return ww.gz.Write(data)
}

func (ww *WrappedWriter) WriteHeader(statuscode int) {
	ww.rw.WriteHeader(statuscode)
}
