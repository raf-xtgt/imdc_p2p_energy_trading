package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/joho/godotenv"
)

func run() error {
	mux := makeMuxRouter()
	httpAddr := os.Getenv("PORT")
	log.Println("Listening on ", os.Getenv("PORT"))
	s := &http.Server{
		Addr:           ":" + httpAddr,
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	if err := s.ListenAndServe(); err != nil {
		return err
	}

	return nil
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
	readjson()
	go func() {
		if Blockchain == nil {
			t := time.Now()
			genesisTransactionData := TransactionData{}
			genesisTransactionData = TransactionData{t.String(), 0}
			spew.Dump(genesisTransactionData)
			pending = append(pending, genesisTransactionData)

			genesisBlock := Block{}
			genesisBlock = Block{0, pending, calculateHash(genesisBlock), "", difficulty, ""}
			spew.Dump(genesisBlock)
			Blockchain = append(Blockchain, genesisBlock)
			clearPending()
		}
	}()
	log.Fatal(run())
}

func readjson() {
	jsonFile, err := os.Open("blockchain.json")
	if err != nil {
		log.Fatal(err)
	}
	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &Blockchain)
}

func writejson() {
	jsonFile, _ := json.MarshalIndent(Blockchain, "", " ")

	_ = ioutil.WriteFile("blockchain.json", jsonFile, 0644)
}

func clearPending() {
	var emptyPending []TransactionData
	pending = emptyPending
}
