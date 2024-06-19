package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const (
	username = "airflow"
	password = "airflow"
	hostname = "127.0.0.1:3306"
	dbname   = "airflow"
)

type CryptoPrice struct {
	ID        uint   `gorm:"primaryKey" json:"id"`
	Symbol    string `json:"Symbol"`
	Name      string `json:"name"`
	Price_usd string `json:"ราคา"`
}

func dsn(dbName string) string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", username, password, hostname, dbName)
}

func getCryptoPrices(w http.ResponseWriter, r *http.Request) {
	db, err := gorm.Open(mysql.Open(dsn(dbname)), &gorm.Config{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var cryptoPrices []CryptoPrice
	if err := db.Find(&cryptoPrices).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cryptoPrices)
}

func createCryptoPrice(w http.ResponseWriter, r *http.Request) {
	db, err := gorm.Open(mysql.Open(dsn(dbname)), &gorm.Config{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var cryptoPrice CryptoPrice
	if err := json.NewDecoder(r.Body).Decode(&cryptoPrice); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := db.Create(&cryptoPrice).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cryptoPrice)
}

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/api/crypto_prices", getCryptoPrices).Methods("GET")
	router.HandleFunc("/api/crypto_prices", createCryptoPrice).Methods("POST")

	log.Fatal(http.ListenAndServe(":8000", router))
}
