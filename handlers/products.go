package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"go-microservice/data"
	"go-microservice/db"
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

	lp, err := db.Get()
	if err != nil {
		http.Error(rw, "Unable to read products", http.StatusInternalServerError)
	}
	marshal, err := json.Marshal(lp)
	if err != nil {
		http.Error(rw, "Unable to marshal json", http.StatusInternalServerError)
	}
	rw.Write(marshal)
}

func (p *Products) AddProduct(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle POST Products")
	prod := r.Context().Value(KeyProduct{}).(data.Product)
	err := db.Insert(&prod)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
	}
	rw.Write([]byte("The product added"))
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

	err = db.Update(&prod)
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
	err = db.Delete(id)
	if err != nil || err == data.ErrProductNotFound {
		http.Error(rw, "Product Not Found", http.StatusNotFound)
		return
	}
	rw.Write([]byte("The product removed"))

}
func (p *Products) GetProductById(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(rw, "Unable to convert the id", http.StatusBadRequest)
		return
	}
	p.l.Println("Handle GetByID Product", id)

	productById, err := db.FindById(id)

	if productById != nil {
		err = productById.ToJSON(rw)
		if err != nil {
			http.Error(rw, "Unable to marshal json", http.StatusInternalServerError)
		}
	} else {
		http.Error(rw, "Product Not Found", http.StatusNotFound)
		return
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

	productByName, err := db.FindByName(name)
	if err == data.ErrProductNotFound {
		http.Error(rw, "Product Not Found", http.StatusNotFound)
		return
	}
	if productByName == nil {
		http.Error(rw, "The product does not exist", http.StatusNotFound)
	} else {
		err = productByName.ToJSON(rw)
		if err != nil {
			http.Error(rw, "Unable to marshal json", http.StatusInternalServerError)
		}
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
