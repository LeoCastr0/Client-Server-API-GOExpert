package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var currenciesDatabase = "currencies.db"

type Currencies struct {
	Currencies ExchangeData `json:"USDBRL"`
}

type ExchangeData struct {
	Code       string `json:"code"`
	Codein     string `json:"codein"`
	Name       string `json:"name"`
	Low        string `json:"low"`
	High       string `json:"high"`
	VarBid     string `json:"varBid"`
	PctChange  string `json:"pctChange"`
	Bid        string `json:"bid"`
	Ask        string `json:"ask"`
	Timestamp  string `json:"timeStamp"`
	CreateDate string `json:"create_date"`
}

type Currency struct {
	ID    int `gorm:"primaryKey"`
	Value float64
	gorm.Model
}

func main() {
	http.HandleFunc("/cotacao", ExchangeRateHandler)
	http.ListenAndServe(":8080", nil)
}

func ExchangeRateHandler(w http.ResponseWriter, r *http.Request) {
	path := "/Users/leonardocastro/Desktop/pos_go/Client-Server-API/server/cotacao.txt"
	if r.URL.Path != "/cotacao" {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	ctx := r.Context()
	select {
	case <-time.After(200 * time.Millisecond):
		log.Println("Request successfully processed")
	case <-ctx.Done():
		log.Println("Request cancelled by client")
		w.WriteHeader(http.StatusInternalServerError)
	}
	log.Println("Request started")
	defer log.Println("Request ended")

	exchangeValue, err := ExchangeRate()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	floatValue, err := strconv.ParseFloat(exchangeValue.Bid, 64)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = AddCurrency(floatValue)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = CreateFile(path)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = SaveFile(exchangeValue.Bid)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	result, err := json.Marshal(exchangeValue.Bid)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(result)
}

func ExchangeRate() (*ExchangeData, error) {
	resp, err := http.Get("https://economia.awesomeapi.com.br/json/last/USD-BRL")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var currencies Currencies
	err = json.Unmarshal(body, &currencies)
	if err != nil {
		return nil, err
	}
	return &currencies.Currencies, nil
}

func SaveFile(currencyValue string) error {
	f, err := os.OpenFile("./cotacao.txt", os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.WriteString(fmt.Sprintf("Dolar: %v\n", currencyValue))
	if err != nil {
		fmt.Printf("Dolar: %v", currencyValue)
		return err
	}
	return nil
}

func CreateFile(path string) error {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		f, err := os.Create("cotacao.txt")
		if err != nil {
			return err
		}
		f.Close()
		return nil
	}
	return nil
}

func initDatabase() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(currenciesDatabase), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	err = db.AutoMigrate(&Currency{})
	if err != nil {
		return nil, err
	}
	return db, nil
}

func AddCurrency(value float64) error {
	db, err := initDatabase()
	if err != nil {
		return err
	}
	db.Create(&Currency{
		Value: value,
	})
	return nil
}
