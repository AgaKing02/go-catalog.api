package main

import (
	"context"
	"github.com/go-openapi/runtime/middleware"
	"github.com/gorilla/mux"
	"go-microservice/db"
	"go-microservice/handlers"
	"gopkg.in/mgo.v2"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	session, err := mgo.Dial("mongodb://localhost:27017")
	if err != nil {
		panic(err)
	}
	defer session.Close()

	db.Database.Context = session

	l := log.New(os.Stdout, "product-api", log.LstdFlags)
	pr := handlers.NewProduct(l)

	sm := mux.NewRouter()
	getRouter := sm.Methods(http.MethodGet).Subrouter()
	putRouter := sm.Methods(http.MethodPut).Subrouter()
	postRouter := sm.Methods(http.MethodPost).Subrouter()
	deleteRouter := sm.Methods(http.MethodDelete).Subrouter()

	putRouter.HandleFunc("/{id:[0-9]+}", pr.UpdateProduct)
	putRouter.Use(pr.MiddlewareProductValidation)

	deleteRouter.HandleFunc("/{id:[0-9]+}", pr.DeleteProduct)

	getRouter.HandleFunc("/", pr.GetProducts)
	getRouter.HandleFunc("/{id:[0-9]+}", pr.GetProductById)

	getRouter.HandleFunc("/name/{name:[A-Za-z]+}", pr.GetProductByName)

	postRouter.HandleFunc("/", pr.AddProduct)
	postRouter.Use(pr.MiddlewareProductValidation)

	opts := middleware.RedocOpts{SpecURL: "/swagger.yaml"}
	sh := middleware.Redoc(opts, nil)

	getRouter.Handle("/docs", sh)

	getRouter.Handle("/swagger.yaml", http.FileServer(http.Dir("./")))

	s := &http.Server{
		Addr:         ":9090",
		Handler:      sm,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}

	go func() {
		err := s.ListenAndServe()
		if err != nil {
			l.Fatal(err)
		}
	}()
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)
	sig := <-sigChan

	l.Println("Received terminate , graceful shutdown", sig)
	tc, _ := context.WithTimeout(context.Background(), 30*time.Second)

	s.Shutdown(tc)

}
