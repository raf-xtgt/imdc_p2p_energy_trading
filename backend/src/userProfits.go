package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func getUserIncomeData(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Methods", "GET,POST,OPTIONS,DELETE,PUT")
	var userId string
	var response string

	// get the data from the request body
	decoder := json.NewDecoder(r.Body)
	// place the user data inside newUser
	if err := decoder.Decode(&userId); err != nil {
		fmt.Println("Failed adding new forecast data", err)
		respondWithJSON(w, r, http.StatusBadRequest, r.Body)
		return
	}
	defer r.Body.Close()
	fmt.Println("This is the user id for getting income data", userId)
	// to prevent backend to timeout
	mongoparams.ctx, mongoparams.cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer mongoparams.cancel()
	respondWithJSON(w, r, http.StatusCreated, response)

}
