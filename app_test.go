package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

var a App

func TestMain(m *testing.M) {

	err := a.Initialize(Dbuser, "kassim.TEstso@me34", "test")
	if err != nil {
		log.Fatal("an error occured while initilizing the database")
	}
	createTable()
	m.Run()
}

func createTable() {
	createTableQuery := `CREATE TABLE IF NOT EXISTS products (
		id INT NOT NULL AUTO_INCREMENT,
		name VARCHAR(255) NOT NULL,
		quantity INT,
		price FLOAT(10,7),
		PRIMARY KEY (id)
	);`

	_, err := a.DB.Exec(createTableQuery)
	if err != nil {
		log.Fatal(err)
	}
}

func clearTable() {
	a.DB.Exec("delete from products")
	a.DB.Exec("alter table products AUTO_INCREMENT=1")
	log.Println("cleared products table")
}

func addProduct(name string, quantity int, price float64) {
	query := "insert into products (name,quantity,price) values(?,?,?)"
	_, err := a.DB.Exec(query, name, quantity, price)
	if err != nil {
		log.Println(err)
	}
}
func TestGetProduct(t *testing.T) {
	clearTable()
	addProduct("keyboard", 100, 200.00)
	request, _ := http.NewRequest("GET", "/product/1", nil)

	response := sendRequest(request)
	checkStatusCode(t, http.StatusOK, response.Code)

}

func checkStatusCode(t *testing.T, expectedStatusCode int, actualStatusCode int) {
	if expectedStatusCode != actualStatusCode {
		t.Errorf("Expected status code %v, but got %v", expectedStatusCode, actualStatusCode)
	}
}
func sendRequest(request *http.Request) *httptest.ResponseRecorder {
	recorder := httptest.NewRecorder()
	a.Router.ServeHTTP(recorder, request)

	return recorder

}

func TestCreateProduct(t *testing.T) {
	clearTable()
	product := []byte(`{"name":"chair","quantity":20,"price":70.2}`)
	request, _ := http.NewRequest("POST", "/product", bytes.NewBuffer(product))
	request.Header.Set("Content-Type", "application/json")
	response := sendRequest(request)
	checkStatusCode(t, http.StatusCreated, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)
	if m["name"] != "chair" {
		t.Errorf("Expected %v, but got %v", "chair", m["name"])
	}

	if m["quantity"] != 20.0 {
		t.Errorf("Expected %v, but got %v ", 20.0, m["quantity"])
	}
	//fmt.Printf("%T, %v", m["price"], m["price"])
	if m["price"] != 70.2 {
		t.Errorf("Expected %v, but got %v", 70.2, m["price"])
	}

}

func TestDeleteProduct(t *testing.T) {
	// delete all products from the database and reset auto_incrementer
	clearTable()
	// add new product to the databse
	addProduct("connector", 80, 140)
	// get the newly added product
	request, _ := http.NewRequest("GET", "/product/1", nil)
	response := sendRequest(request)
	checkStatusCode(t, http.StatusOK, response.Code)
	/// delete the product
	request, _ = http.NewRequest("DELETE", "/product/1", nil)
	sendRequest(request)

	// check if the product exists again

	request, _ = http.NewRequest("GET", "/product/1", nil)
	response = sendRequest(request)
	checkStatusCode(t, http.StatusNotFound, response.Code)

}

func TestUpdateProduct(t *testing.T) {
	// clear the database
	clearTable()
	// add new product to the databse
	addProduct("Hp Laptop", 30, 560)

	// get the product
	request, _ := http.NewRequest("GET", "/product/1", nil)
	response := sendRequest(request)
	checkStatusCode(t, http.StatusOK, response.Code)
	var oldValue map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &oldValue)

	//update the databse
	var product = []byte(`{"name":"lenova laptop","quantity":29,"price":345}`)
	request, _ = http.NewRequest("PUT", "/product/1", bytes.NewBuffer(product))
	request.Header.Set("Content-Type", "application/json")
	response = sendRequest(request)
	var newValue map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &newValue)

	if oldValue["id"] != newValue["id"] {
		t.Errorf("Expected id %v, but got %v", newValue["id"], oldValue["id"])
	}

	if oldValue["quantity"] == newValue["quanity"] {
		t.Errorf("quanity not updated")
	}

	if oldValue["price"] == newValue["quantity"] {
		t.Errorf("price was not updated")
	}

}
