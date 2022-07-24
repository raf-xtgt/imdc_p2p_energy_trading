package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

/**
Function to create a database entry for user balance
ASSUMPTION 1: The energy balance comes from the smart metre
ASSUMPTION 2: We assume already in sync with some financial institution
**/
func createUserAccount(userId string) {
	// to prevent backend to timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var balance AccountBalance
	balance.UserId = userId
	// add dummy data
	balance.FiatBalance = 20000
	balance.EnergyBalance = 800

	//write the balance in the database info to the users collection
	writeData, err := db.UserAccBalance.InsertOne(ctx, balance)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("User balance added", writeData.InsertedID)
	return

}

// Get the orders of the user....the user is the buyer

func getUserBuyRequests(w http.ResponseWriter, r *http.Request) {

	var userId string
	var response TransactionsResponse
	w.Header().Add("Access-Control-Allow-Methods", "GET,POST,OPTIONS,DELETE,PUT")

	// get the data from json body
	decoder := json.NewDecoder(r.Body)
	// place the user data inside newRequest
	if err := decoder.Decode(&userId); err != nil {
		fmt.Println("Failed retrieving the user id from frontend: Line --> 48", err)
		respondWithJSON(w, r, http.StatusBadRequest, r.Body)
		return
	}

	defer r.Body.Close()

	// to prevent backend to timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// get the transactions that have not been verified yet
	cursor, err := db.Transactions.Find(ctx, bson.M{"buyerId": userId})
	if err != nil {
		log.Fatal(err)
	}

	var transactions []Transaction
	if err = cursor.All(ctx, &transactions); err != nil {
		log.Fatal(err)
	}

	response.Transactions = transactions
	respondWithJSON(w, r, http.StatusCreated, response)

}

// Send all the transactions to the frontend
func getUserSellRequests(w http.ResponseWriter, r *http.Request) {

	var userId string
	var response TransactionsResponse
	w.Header().Add("Access-Control-Allow-Methods", "GET,POST,OPTIONS,DELETE,PUT")

	// get the data from json body
	decoder := json.NewDecoder(r.Body)
	// place the user data inside newRequest
	if err := decoder.Decode(&userId); err != nil {
		fmt.Println("Failed retrieving the user id from frontend: Line --> 48", err)
		respondWithJSON(w, r, http.StatusBadRequest, r.Body)
		return
	}

	defer r.Body.Close()

	// to prevent backend to timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// get the transactions that have not been verified yet
	cursor, err := db.Transactions.Find(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}

	var transactions []Transaction
	if err = cursor.All(ctx, &transactions); err != nil {
		log.Fatal(err)
	}

	response.Transactions = transactions
	respondWithJSON(w, r, http.StatusCreated, response)

}

/**

{
_id:6264d05456771e5a219a0a56

buyerId: "623722a33811066ba71aed41"

buyerPayable: 34.480434782608704

buyerEnReceivableFromAuction: 22.739543331907836

buyerEnReceivableFromTNB: 127.26045666809216

auctionBids: <Array>

TNBReceivable: 29.253305843661018

verified: true

chained: false

tId: "6264d05456771e5a219a0a56"

checks: 2

date: "20-04-2022"

}
*/
