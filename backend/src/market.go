package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
	//crypto library to hash the data
	// to create the secret(private) key for jwt token
)

// function to store a new buy request in the database
func addBuyRequest(w http.ResponseWriter, r *http.Request) {
	var newRequest BuyRequest
	//w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Methods", "GET,POST,OPTIONS,DELETE,PUT")

	// get the data from json body
	decoder := json.NewDecoder(r.Body)
	// place the user data inside newRequest
	if err := decoder.Decode(&newRequest); err != nil {
		fmt.Println("Failed adding a new request", err)
		respondWithJSON(w, r, http.StatusBadRequest, r.Body)
		return
	}
	defer r.Body.Close()

	// to prevent backend to timeout
	mongoparams.ctx, mongoparams.cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer mongoparams.cancel()

	//write user info to the users collection
	writeRequest, err := db.EnergyBuyRequests.InsertOne(mongoparams.ctx, newRequest)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("New energy request added", writeRequest.InsertedID)

	respondWithJSON(w, r, http.StatusCreated, newRequest)
	//respondWithJSON(w, r, http.StatusCreated, NewUser)
	return
}
