package main

import (
	"encoding/json"
	"log"
	"net/http"
	"fmt"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/lib/pq"
    "github.com/rs/cors"

	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type Resource struct {
	gorm.Model

	Link        string
	Name        string
	Author      string
	Description string
	Tags        pq.StringArray	 `gorm:"type:varchar(64)[]"`
}

var db *gorm.DB
var err error

func main() {
	router := mux.NewRouter()
	db, err = gorm.Open("postgres", "host=localhost port=5432 user=postgres dbname=calhounio_demo password=postgres")
	if err != nil {
		panic("failed to connect database")
	}

	defer db.Close()

	db.AutoMigrate(&Resource{})
	fmt.Println("DB set")

	router.HandleFunc("/resources", GetResources).Methods("GET")
	router.HandleFunc("/resources/{id}", GetResource).Methods("GET")
	router.HandleFunc("/resources", CreateResource).Methods("POST")
	router.HandleFunc("/resources/{id}", DeleteResource).Methods("DELETE")

    handler := cors.Default().Handler(router)
	
	log.Fatal(http.ListenAndServe(":8000", handler))
}

func GetResources(w http.ResponseWriter, r *http.Request) {
	var resources []Resource
	db.Find(&resources)
	json.NewEncoder(w).Encode(&resources)
}

func GetResource(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var resource Resource
	db.First(&resource, params["id"])
	json.NewEncoder(w).Encode(&resource)
}

func CreateResource(w http.ResponseWriter, r *http.Request) {
	var resource Resource
	json.NewDecoder(r.Body).Decode(&resource)
	db.Create(&resource)
	json.NewEncoder(w).Encode(&resource)
}

func DeleteResource(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var resource Resource
	db.First(&resource, params["id"])
	db.Delete(&resource)
	fmt.Println(params["id"] + " Deleted")

	var resources []Resource
	db.Find(&resources)
	json.NewEncoder(w).Encode(&resources)
}