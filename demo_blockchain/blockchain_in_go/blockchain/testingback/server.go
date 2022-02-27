package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/mux"
)

var mutex = &sync.Mutex{}

func makeMuxRouter() http.Handler {
	muxRouter := mux.NewRouter()
	muxRouter.HandleFunc("/", handleGetBlockchain).Methods("GET")
	muxRouter.HandleFunc("/WriteBlock", handleWriteBlock)
	muxRouter.HandleFunc("/print", Validate)
	muxRouter.HandleFunc("/AddUser", Adduser)
	muxRouter.HandleFunc("/AddLeader", AddLeader)
	muxRouter.HandleFunc("/AddValidator", AddValidator)
	muxRouter.HandleFunc("/StartServer", StartTCPServer)
	muxRouter.HandleFunc("/SendSomething", ValidateBlockchain)
	return muxRouter
}

// Obtain from mycoralhealth website
func handleGetBlockchain(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Methods", "GET,POST,OPTIONS,DELETE,PUT")
	bytes, err := json.MarshalIndent(Blockchain, "", "  ")
	bytes2, err2 := json.MarshalIndent(pending, "", "  ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err2 != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	io.WriteString(w, string(bytes)+string(bytes2))
}

// Generate a Block from list of Transaction Data and send it to the other validators
func Validate(w http.ResponseWriter, r *http.Request) {

	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Methods", "GET,POST,OPTIONS,DELETE,PUT")

	newBlock := generateBlock(Blockchain[len(Blockchain)-1])

	// Send the newBlock to the other validators
	SendBlock(newBlock)
	writejson()

	// Clear and reset the pending array
	clearPending()

}

// Create a transaction data block from the HTTP POST request
func handleWriteBlock(w http.ResponseWriter, r *http.Request) {
	var data2 AltTransactionData
	var newTransactionData TransactionData

	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Methods", "GET,POST,OPTIONS,DELETE,PUT")
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&data2); err != nil {
		respondWithJSON(w, r, http.StatusBadRequest, r.Body)
		return
	}
	defer r.Body.Close()

	//ensure atomicity when creating new block
	mutex.Lock()
	newTransactionData.Buyer, _ = strconv.Atoi(data2.Buyer)
	newTransactionData.Seller, _ = strconv.Atoi(data2.Seller)
	newTransactionData.Money, _ = strconv.Atoi(data2.Money)
	newTransactionData.Energy, _ = strconv.Atoi(data2.Energy)
	t := time.Now()
	newTransactionData.Timestamp = t.String()

	// Validate the Transaction data
	TxValidity := validateTransaction(newTransactionData)
	if TxValidity {
		newPending := append(pending, newTransactionData)
		pending = newPending
	} else {
		fmt.Println("Transaction error")
	}
	mutex.Unlock()

}

// Obtain from mycoralhealth website, simply return the HTTP status and data to the front end
func respondWithJSON(w http.ResponseWriter, r *http.Request, code int, payload interface{}) {
	response, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("HTTP 500: Internal Server Error"))
		return
	}
	w.WriteHeader(code)
	w.Write(response)
}

// Add new user, leader, validator
func Adduser(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Methods", "GET,POST,OPTIONS,DELETE,PUT")
	CreateNewUser()
}

func AddLeader(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Methods", "GET,POST,OPTIONS,DELETE,PUT")
	CreateLeader()
}

func AddValidator(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Methods", "GET,POST,OPTIONS,DELETE,PUT")
	CreateValidator()
}

// manually start the TCP port 3000
func StartTCPServer(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Methods", "GET,POST,OPTIONS,DELETE,PUT")

	fmt.Println("StartLeaderServer")
	StartLeaderServer()
}

// manually check the blockchain
func ValidateBlockchain(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Methods", "GET,POST,OPTIONS,DELETE,PUT")

	checkBlockchainValid()
}
