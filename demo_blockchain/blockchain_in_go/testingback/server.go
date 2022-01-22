package main

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"sync"

	"github.com/davecgh/go-spew/spew"
	"github.com/gorilla/mux"
)

var mutex = &sync.Mutex{}

func makeMuxRouter() http.Handler {
	muxRouter := mux.NewRouter()
	muxRouter.HandleFunc("/", handleGetBlockchain).Methods("GET")
	muxRouter.HandleFunc("/WriteBlock", handleWriteBlock)
	muxRouter.HandleFunc("/print", printTime)
	return muxRouter
}

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

func printTime(w http.ResponseWriter, r *http.Request) {

	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Methods", "GET,POST,OPTIONS,DELETE,PUT")
	newBlock := generateBlock(Blockchain[len(Blockchain)-1])
	if isBlockValid(newBlock, Blockchain[len(Blockchain)-1]) {
		newBlockchain := append(Blockchain, newBlock)
		replaceChain(newBlockchain)
		writejson()
		clearPending()
	}
}

type Message struct {
	BPM int
}

func handleWriteBlock(w http.ResponseWriter, r *http.Request) {
	var m Message
	var BPM1 string

	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Methods", "GET,POST,OPTIONS,DELETE,PUT")

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&BPM1); err != nil {
		respondWithJSON(w, r, http.StatusBadRequest, r.Body)
		return
	}
	defer r.Body.Close()

	BPM2, err := strconv.Atoi(BPM1)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	m.BPM = BPM2
	//ensure atomicity when creating new block
	mutex.Lock()
	newTransactionData := generateTransactionData(m.BPM)
	mutex.Unlock()
	if err != nil {
		respondWithJSON(w, r, http.StatusInternalServerError, m)
		return
	}
	newBlockchain := append(pending, newTransactionData)
	replacePendingChain(newBlockchain)
	spew.Dump(Blockchain)

	respondWithJSON(w, r, http.StatusCreated, newTransactionData)

}

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
