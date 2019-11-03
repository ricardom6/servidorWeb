package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func homeLink(W http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(W, "Welcome Bank!")
}

type account struct {
	ID     string    `json:"ID"`
	Saldo  float32   `json:"Saldo"`
	Data   time.Time `json:"Data"`
	Status bool      `json:"Status"`
}

type allaccounts []account

var accounts = allaccounts{
	{
		ID:     "1",
		Saldo:  123.45,
		Data:   time.Now(),
		Status: true,
	},
}

func createaccount(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	var newaccount account

	reqBody, err := ioutil.ReadAll(request.Body)
	if err != nil {
		fmt.Fprintf(response, "Informe dados da conta para cadastro.")
	}
	json.Unmarshal(reqBody, &newaccount)
	newaccount.Data = time.Now()
	accounts = append(accounts, newaccount)
	response.WriteHeader(http.StatusCreated)
	json.NewEncoder(response).Encode(newaccount)
}
func getaccount(response http.ResponseWriter, request *http.Request) {
	accountID := mux.Vars(request)["id"]

	for _, singleaccount := range accounts {
		if singleaccount.ID == accountID {
			json.NewEncoder(response).Encode(singleaccount)
		}
	}
}

func getAllaccounts(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	json.NewEncoder(response).Encode(accounts)
}

func creditaccount(response http.ResponseWriter, request *http.Request) {
	accountID := mux.Vars(request)["id"]

	var credit account
	reqBody, err := ioutil.ReadAll(request.Body)
	if err != nil {
		fmt.Fprintf(response, "Informe os dados da conta.")
	}

	json.Unmarshal(reqBody, &credit)

	for i, singleaccount := range accounts {
		if singleaccount.ID == accountID {
			if singleaccount.Status == true {
				singleaccount.Saldo += credit.Saldo
				accounts = append(accounts[:i], singleaccount)
				json.NewEncoder(response).Encode(singleaccount)
				response.WriteHeader(http.StatusNoContent)
			}
		}
	}
}

func debitaccount(response http.ResponseWriter, request *http.Request) {
	accountID := mux.Vars(request)["id"]

	var credit account
	reqBody, err := ioutil.ReadAll(request.Body)
	if err != nil {
		fmt.Fprintf(response, "Informe os dados da conta.")
	}

	json.Unmarshal(reqBody, &credit)

	for i, singleaccount := range accounts {
		if singleaccount.ID == accountID {
			singleaccount.Saldo -= credit.Saldo
			accounts = append(accounts[:i], singleaccount)
			response.WriteHeader(http.StatusNoContent)
			json.NewEncoder(response).Encode(singleaccount)
		}
	}
}

func blockaccount(response http.ResponseWriter, request *http.Request) {
	accountID := mux.Vars(request)["id"]

	var block account
	reqBody, err := ioutil.ReadAll(request.Body)
	if err != nil {
		fmt.Fprintf(response, "Informe os dados da conta.")
	}

	json.Unmarshal(reqBody, &block)

	for i, singleaccount := range accounts {
		if singleaccount.ID == accountID {
			singleaccount.Status = block.Status
			accounts = append(accounts[:i], singleaccount)
			response.WriteHeader(http.StatusNoContent)
			json.NewEncoder(response).Encode(singleaccount)
		}
	}
}

func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", homeLink)
	router.HandleFunc("/api/v1/account", createaccount).Methods("POST")
	router.HandleFunc("/api/v1/account/{id}", getaccount).Methods("GET")
	router.HandleFunc("/api/v1/account/all", getAllaccounts).Methods("GET")
	router.HandleFunc("/api/v1/account/{id}/block", blockaccount).Methods("PATCH")
	router.HandleFunc("/api/v1/account/{id}/credit", creditaccount).Methods("PATCH")
	router.HandleFunc("/api/v1/account/{id}/debit", debitaccount).Methods("PATCH")
	fmt.Printf("Servidor disponivel!")
	log.Fatal(http.ListenAndServe(":8080", router))

}
