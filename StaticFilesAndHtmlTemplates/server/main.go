package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/mux"
)

type DataToCommonWindow struct {
	Title            string
	ChoiceCountry    string
	ListCountry      []Country
	DataFromResponse []CurrencyData
}
type Country struct {
	Name string
	Code string
}
type CurrencyData struct {
	Amount float64            `json:"amount"`
	Base   string             `json:"base"`
	Date   string             `json:"date"`
	Rates  map[string]float64 `json:"rates"`
}

var datachan = make(chan []CurrencyData)
var CurrencyDataList []CurrencyData
var ListToGetter = []Country{
	{Name: "Доллар США", Code: "USD"}, {Name: "Евро", Code: "EUR"},
	{Name: "Йена", Code: "JPY"}, {Name: "Юань", Code: "CNY"}, {Name: "Стерлинг", Code: "GBP"}, {Name: "Франк", Code: "CHF"},
}
var ResponseToServer []CurrencyData

func GetValueFromApi() ([]CurrencyData, error) {
	for _, cou := range ListToGetter {
		for _, m := range ListToGetter {
			if cou != m {
				var newdata CurrencyData
				newurl := "https://api.frankfurter.app/latest?from=" + cou.Code + "&to=" + m.Code
				response, err := http.Get(newurl)
				if err != nil {
					log.Println(err)
					return nil, err
				}
				defer response.Body.Close()
				err = json.NewDecoder(response.Body).Decode(&newdata)
				if err != nil {
					log.Println(err)
					return nil, err
				}
				CurrencyDataList = append(CurrencyDataList, newdata)
			}
		}
	}
	return CurrencyDataList, nil
}
func TimerDoValueFromApi() {
	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()
	for range ticker.C {
		data, err := GetValueFromApi()
		if err != nil {
			log.Println(err)
			continue
		}
		datachan <- data
	}
}
func updateCurrencyData() {
	for data := range datachan {
		ResponseToServer = data
		log.Println("Currency data updated")
	}
}
func CommonWindowGet(we http.ResponseWriter, r *http.Request) {
	dt := DataToCommonWindow{Title: "Курс валют", ChoiceCountry: "Выберите валюту:", ListCountry: make([]Country, 0)}
	dt.ListCountry = append(dt.ListCountry, Country{Name: "Доллар США", Code: "USD"}, Country{Name: "Евро", Code: "EUR"},
		Country{Name: "Йена", Code: "JPY"}, Country{Name: "Юань", Code: "CNY"}, Country{Name: "Стерлинг", Code: "GBP"}, Country{Name: "Франк", Code: "CHF"},
	)
	dt.DataFromResponse = ResponseToServer
	tmp, err := template.ParseFiles("static/porka.html")
	if err != nil {
		log.Println(err)
		return
	}
	err = tmp.Execute(we, dt)
	if err != nil {
		log.Println(err)
		return
	}
}

func main() {
	go TimerDoValueFromApi()
	go updateCurrencyData()
	initialData, err := GetValueFromApi()
	if err != nil {
		log.Println("Error fetching initial data:", err)
	} else {
		datachan <- initialData
	}
	router := mux.NewRouter()
	server := http.Server{
		Addr:    ":7723",
		Handler: router,
	}
	router.HandleFunc("/", CommonWindowGet).Methods("GET")
	log.Fatal(server.ListenAndServe())
	r := gin.Default()
	r.POST("", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"msg": "ftd",
		})
	})
	r.Run(":7789")
}
