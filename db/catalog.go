package db

import (
	"fmt"
	"go-microservice/data"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type ProductDb struct {
	Context *mgo.Session
}

var Database = ProductDb{}

func Insert(p *data.Product) error {
	productCollection := Database.Context.DB("productDb").C("products")
	err := productCollection.Insert(p)
	if err != nil {
		return err
	}
	fmt.Println("Successfully Inserted")
	return nil
}
func Delete(id int) error {
	productCollection := Database.Context.DB("productDb").C("products")
	filter := bson.M{"id": id}
	err := productCollection.Remove(filter)
	if err != nil {
		return err
	}
	fmt.Println("Successfully Deleted")

	return nil
}

func FindById(id int) (*data.Product, error) {
	query := bson.M{"id": id}
	productCollection := Database.Context.DB("productDb").C("products")
	prod := &data.Product{}
	err := productCollection.Find(query).One(prod)
	if err != nil {
		return nil, err
	}
	return prod, nil
}

func FindByName(name string) (*data.Product, error) {
	query := bson.M{"name": name}
	productCollection := Database.Context.DB("productDb").C("products")
	prod := &data.Product{}
	err := productCollection.Find(query).One(prod)
	if err != nil {
		return nil, err
	}
	fmt.Println(prod)
	return prod, nil
}

func FindByCategory(category string) (*[]data.Product, error) {
	query := bson.M{"category": category}
	productCollection := Database.Context.DB("productDb").C("products")
	products := &[]data.Product{}
	err := productCollection.Find(query).All(&products)
	if err != nil {
		return nil, err
	}
	return products, nil
}

func Get() (*[]data.Product, error) {
	query := bson.M{}
	productCollection := Database.Context.DB("productDb").C("products")
	productArray := &[]data.Product{}
	err := productCollection.Find(query).All(&productArray)
	if err != nil {
		return nil, err
	}
	return productArray, nil
}

func Update(p *data.Product) error {
	productCollection := Database.Context.DB("productDb").C("products")
	filter := bson.M{"id": p.ID}
	uProd := bson.M{"$set": bson.M{"name": p.Name, "description": p.Description, "sku": p.SKU, "createdon": p.CreatedOn, "updatedon": p.UpdatedOn, "deletedon": p.DeletedOn}}
	err := productCollection.Update(filter, uProd)
	if err != nil {
		return err
	}
	fmt.Println("Successfully Updated")

	return nil
}
