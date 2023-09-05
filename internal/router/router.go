package router

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"regexp"
	"strconv"

	"example.com/apiserver/internal/db"
	"github.com/gorilla/mux"
)

var conn *sql.DB

func getProducts(w http.ResponseWriter, r *http.Request) {
	log.Println("Endpoint hit: returnALlProducts")
	products, err := db.GetProducts(conn)
	if err != nil {
		sendError(w, http.StatusInternalServerError, err.Error())
		return
	}
	sendResponse(w, http.StatusOK, products)
}

func getProduct(w http.ResponseWriter, r *http.Request) {
	log.Println("Endpoint hit: getProduct")
	var id string
	regexCompiler := regexp.MustCompile(`/product/(.*)`)
	match := regexCompiler.FindStringSubmatch(r.URL.Path)
	if len(match) > 1 {
		id = match[1]
	}
	productId, err := strconv.Atoi(id)
	if err != nil {
		sendError(w, http.StatusBadRequest, "Id should be a number")
		return
	}
	product, err := db.GetProduct(conn, productId)
	if err != nil {
		sendError(w, http.StatusInternalServerError, err.Error())
		return
	}
	sendResponse(w, http.StatusOK, product)
}

func createProduct(w http.ResponseWriter, r *http.Request) {
	log.Println("Endponthit: createProduct")
	var p db.Product
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		sendError(w, http.StatusBadRequest, err.Error())
		return
	}
	err = db.CreateProduct(conn, &p)
	if err != nil {
		sendError(w, http.StatusInternalServerError, err.Error())
		return
	}
	sendResponse(w, http.StatusCreated, p)
}

func updateProduct(w http.ResponseWriter, r *http.Request) {
	log.Println("Endpointhit: updateProduct")
	var id string
	regexCompiler := regexp.MustCompile(`/product/(.*)`)
	match := regexCompiler.FindStringSubmatch(r.URL.Path)
	if len(match) > 1 {
		id = match[1]
	}
	productId, err := strconv.Atoi(id)
	if err != nil {
		sendError(w, http.StatusBadRequest, "Id should be a number")
	}
	product, err := db.GetProduct(conn, productId)
	if err != nil {
		sendError(w, http.StatusInternalServerError, err.Error())
		return
	}
	err = json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		sendError(w, http.StatusBadRequest, err.Error())
		return
	}
	err = db.UpdateProduct(conn, product)
	if err != nil {
		sendError(w, http.StatusInternalServerError, err.Error())
		return
	}
	sendResponse(w, http.StatusOK, product)
}

func deleteProduct(w http.ResponseWriter, r *http.Request) {
	log.Println("Endpointhit: deleteProduct")
	var id string
	regexCompiler := regexp.MustCompile(`/product/(.*)`)
	match := regexCompiler.FindStringSubmatch(r.URL.Path)
	if len(match) > 1 {
		id = match[1]
	}
	productId, err := strconv.Atoi(id)
	if err != nil {
		sendError(w, http.StatusBadRequest, "Id should be a number")
	}
	err = db.DeleteProduct(conn, productId)
	if err != nil {
		sendError(w, http.StatusInternalServerError, err.Error())
		return
	}
	sendResponse(w, http.StatusOK, map[string]string{"message": "successfully deleted"})
}

func sendResponse(w http.ResponseWriter, statusCode int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(response)
}

func sendError(w http.ResponseWriter, statusCode int, err string) {
	errorMessage := map[string]string{"error": err}
	sendResponse(w, statusCode, errorMessage)
}

func addRoutes(r *mux.Router) {
	r.HandleFunc("/products", getProducts).Methods(http.MethodGet)
	r.HandleFunc("/product/{id}", getProduct).Methods(http.MethodGet)
	r.HandleFunc("/product", createProduct).Methods(http.MethodPost)
	r.HandleFunc("/product/{id}", updateProduct).Methods(http.MethodPut)
	r.HandleFunc("/product/{id}", deleteProduct).Methods(http.MethodDelete)
}

func Initialise(db *sql.DB) *mux.Router {
	conn = db
	router := mux.NewRouter()
	addRoutes(router)
	return router
}
