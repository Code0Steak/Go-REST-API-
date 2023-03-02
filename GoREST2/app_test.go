package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

var a App

func TestMain(m *testing.M) {

	err := a.Initialize(DB_Username, DB_Pass, "test")

	if err != nil {
		log.Fatal("Error occured while initializing the DB.")
	}

	createTable()

	//run all tests
	m.Run()
}

func createTable() {
	query := `CREATE TABLE IF NOT EXISTS products (
		id int NOT NULL AUTO_INCREMENT,
		name varchar(255) NOT NULL,
		quantity int,
		price float(10,7),
		PRIMARY KEY (id)
	);
	`

	_, err := a.DB.Exec(query)
	if err != nil {
		log.Fatal(err)
	}
}

func clearTable() {
	a.DB.Exec("DELETE from products")
	a.DB.Exec("ALTER table products AUTO_INCREMENT=1")
}

func addProduct(name string, quantity int, price float64) {

	query := fmt.Sprintf("INSERT INTO products(name, quantity, price) VALUES(\"%v\",%v,%v)", name, quantity, price)
	_, err := a.DB.Exec(query)

	if err != nil {
		log.Println(err)
	}

}

// Test the 'GET' request for the application
// Steps :- Clear Table, Add new product data, send request, validate/return errors
func TestGetProduct(t *testing.T) {
	clearTable()
	addProduct("Chewing Gum", 30, 9.9)
	request, _ := http.NewRequest("GET", "/product/1", nil) //create new request
	response := sendRequest(request)                        //send request and get the response
	//check status code of the response
	checkStatCode(t, http.StatusOK, response.Code)
}

// create a recorder object and send the request
func sendRequest(r *http.Request) *httptest.ResponseRecorder {
	recorder := httptest.NewRecorder()
	a.Router.ServeHTTP(recorder, r) //recorder is the response writer.
	return recorder
}

// checks status code
func checkStatCode(t *testing.T, expectedStatusCode, actualStatusCode int) {
	if expectedStatusCode != actualStatusCode {
		t.Errorf("Expected status : %v, received : %v", expectedStatusCode, actualStatusCode)
	}
}

//-------------------------------------------------------------------------------------------------------------------------------------------

// Test for 'POST'
func TestCreateProduct(t *testing.T) {
	clearTable()
	var product = []byte(`{"name":"Backpack","quantity":9,"price":500}`)
	req, _ := http.NewRequest("POST", "/product", bytes.NewBuffer(product))
	req.Header.Set("Content-type", "application/json")

	response := sendRequest(req)

	checkStatCode(t, http.StatusCreated, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	//check if the response contains the data that was actually sent.
	if m["name"] != "chair" {
		t.Errorf("Expected : %v, received %v.", "chair", m["name"])
	}

}
