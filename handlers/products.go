package handlers

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"go-microservice/data"
	"log"
	"net/http"
	"strconv"
)

type Products struct {
	l *log.Logger
}
type productResponse struct {
	Body []data.Product
}
type productIDParameterWrapper struct {
	ID int `json:"id"`
}
type productsNoContent struct {
}

func NewProduct(l *log.Logger) *Products {
	return &Products{l}
}

func (p *Products) GetProducts(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle GET Products")

	lp := data.GetProducts()
	err := lp.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to marshal json", http.StatusInternalServerError)
	}
}

func (p *Products) AddProduct(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle POST Products")
	prod := r.Context().Value(KeyProduct{}).(data.Product)
	data.AddProduct(&prod)
}

func (p *Products) UpdateProduct(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(rw, "Unable to convert the id", http.StatusBadRequest)
		return
	}
	p.l.Println("Handle PUT Product", id)
	prod := r.Context().Value(KeyProduct{}).(data.Product)

	err = prod.FromJSON(r.Body)
	if err != nil {
		http.Error(rw, "Unable to unmarshal json", http.StatusBadRequest)
	}
	p.l.Printf("Prod: %#v", prod)

	err = data.UpdateProduct(id, &prod)
	if err == data.ErrProductNotFound {
		http.Error(rw, "Product Not Found", http.StatusNotFound)
		return
	}
}

func (p *Products) DeleteProduct(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(rw, "Unable to convert the id", http.StatusBadRequest)
		return
	}
	p.l.Println("Handle DELETE Product", id)
	//prod :=r.Context().Value(KeyProduct{}).(data.Product)
	//
	//err = prod.FromJSON(r.Body)
	//if err != nil {
	//	http.Error(rw, "Unable to unmarshal json", http.StatusBadRequest)
	//}

	err = data.DeleteProduct(id)
	if err == data.ErrProductNotFound {
		http.Error(rw, "Product Not Found", http.StatusNotFound)
		return
	}
}
func (p *Products) GetProductById(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(rw, "Unable to convert the id", http.StatusBadRequest)
		return
	}
	p.l.Println("Handle GetByID Product", id)
	//prod :=r.Context().Value(KeyProduct{}).(data.Product)
	//
	//err = prod.FromJSON(r.Body)
	//if err != nil {
	//	http.Error(rw, "Unable to unmarshal json", http.StatusBadRequest)
	//}

	productById, err := data.GetProductById(id)
	if err != nil || err == data.ErrProductNotFound {
		http.Error(rw, "Product Not Found", http.StatusNotFound)
		return
	}
	err = productById.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to marshal json", http.StatusInternalServerError)
	}
}

func (p *Products) GetProductByName(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]

	if name == "" {
		http.Error(rw, "Name is null", http.StatusBadRequest)
		return
	}
	p.l.Println("Handle GetByName Product", name)

	productByName, err := data.GetProductByName(name)
	if err != nil || err == data.ErrProductNotFound {
		http.Error(rw, "Product Not Found", http.StatusNotFound)
		return
	}
	err = productByName.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to marshal json", http.StatusInternalServerError)
	}

}

type KeyProduct struct {
}

func (p Products) MiddlewareProductValidation(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		prod := data.Product{}
		err := prod.FromJSON(r.Body)
		if err != nil {
			p.l.Println("Error deserializing the product", err)
			http.Error(rw, "Unable to read product", http.StatusBadRequest)
			return

		}
		err = prod.Validate()
		if err != nil {
			p.l.Println("Error validating the product", err)
			http.Error(rw, fmt.Sprintf("Unable to validate product: %s", err), http.StatusBadRequest)
			return
		}
		ctx := context.WithValue(r.Context(), KeyProduct{}, prod)
		r = r.WithContext(ctx)

		next.ServeHTTP(rw, r)
	})
}
