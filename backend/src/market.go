package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
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
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	//write user info to the users collection\
	loc, _ := time.LoadLocation("UTC")
	now := time.Now().In(loc)
	//requestTime := now
	//fmt.Println("ZONE : ", loc, " Time : ", now) // UTC
	newRequest.ReqTime = now.String()
	newRequest.Auctioned = false
	writeRequest, err := db.EnergyBuyRequests.InsertOne(ctx, newRequest)
	if err != nil {
		log.Fatal(err)
	}

	if addUniqueBuyReqId(newRequest.BuyerId, newRequest.ReqTime, fmt.Sprintf("%v", writeRequest.InsertedID)) {
		fmt.Println("New energy request added", writeRequest.InsertedID)
	}

	respondWithJSON(w, r, http.StatusCreated, newRequest)
	//respondWithJSON(w, r, http.StatusCreated, NewUser)
	return
}

// function to add a unique id to user document
func addUniqueBuyReqId(buyerId string, reqTime string, uniqueId string) bool {

	// slice the id to retain the id part only
	unId := uniqueId[10 : len(uniqueId)-2]
	//uniqueId = unId
	fmt.Println(uniqueId)
	fmt.Println(unId)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// update the document that matches the buyerid and time of order
	_, err := db.EnergyBuyRequests.UpdateOne(
		ctx,
		bson.M{"buyerid": buyerId, "reqtime": reqTime},
		bson.D{
			{"$set", bson.D{{"reqid", unId}}},
		},
	)

	// if the update fails
	if err != nil {
		log.Fatal(err)
		return false
	}
	return true
}

// function to get the energy requests
func getEnergyRequests(w http.ResponseWriter, r *http.Request) {
	//runLa()
	w.Header().Add("Access-Control-Allow-Methods", "GET,POST,OPTIONS,DELETE,PUT")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := db.EnergyBuyRequests.Find(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}

	var allBuyRequests []BuyRequest
	if err = cursor.All(ctx, &allBuyRequests); err != nil {
		log.Fatal(err)
	}
	var response GetBuyEnergyResponse
	response.Requests = allBuyRequests
	respondWithJSON(w, r, http.StatusCreated, response)

	//fmt.Println("Getting all buy requests for market", allBuyRequests)
}

// function to store a new buy request in the database
func addSellRequest(w http.ResponseWriter, r *http.Request) {
	var newRequest SellRequest
	//w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Methods", "GET,POST,OPTIONS,DELETE,PUT")

	// get the data from json body
	decoder := json.NewDecoder(r.Body)
	// place the user data inside newRequest
	if err := decoder.Decode(&newRequest); err != nil {
		fmt.Println("Failed adding a new sell request", err)
		respondWithJSON(w, r, http.StatusBadRequest, r.Body)
		return
	}
	defer r.Body.Close()

	// to prevent backend to timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	//write user info to the users collection\
	loc, _ := time.LoadLocation("UTC")
	now := time.Now().In(loc)
	//requestTime := now
	//fmt.Println("ZONE : ", loc, " Time : ", now) // UTC
	newRequest.ReqTime = now.String()
	writeRequest, err := db.EnergySellRequests.InsertOne(ctx, newRequest)
	if err != nil {
		log.Fatal(err)
	}

	if addUniqueSellReqId(newRequest.SellerId, newRequest.ReqTime, fmt.Sprintf("%v", writeRequest.InsertedID)) {
		fmt.Println("New sell energy request added with id", writeRequest.InsertedID)
	}

	respondWithJSON(w, r, http.StatusCreated, newRequest)
	//respondWithJSON(w, r, http.StatusCreated, NewUser)
	return
}

// function to add a unique id to user document
func addUniqueSellReqId(buyerId string, reqTime string, uniqueId string) bool {

	// slice the id to retain the id part only
	unId := uniqueId[10 : len(uniqueId)-2]
	//uniqueId = unId
	fmt.Println(uniqueId)
	fmt.Println(unId)

	// to prevent backend to timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// update the document that matches the buyerid and time of order
	_, err := db.EnergySellRequests.UpdateOne(
		ctx,
		bson.M{"sellerid": buyerId, "reqtime": reqTime},
		bson.D{
			{"$set", bson.D{{"sellreqid", unId}}},
		},
	)

	// if the update fails
	if err != nil {
		log.Fatal(err)
		return false
	}
	return true
}

// function to close a buy request by updating it
func closeBuyRequest(w http.ResponseWriter, r *http.Request) {
	var reqId string
	w.Header().Add("Access-Control-Allow-Methods", "GET,POST,OPTIONS,DELETE,PUT")

	// get the data from json body
	decoder := json.NewDecoder(r.Body)
	// place the user data inside the string
	if err := decoder.Decode(&reqId); err != nil {
		fmt.Println("Failed adding a new sell request", err)
		respondWithJSON(w, r, http.StatusBadRequest, r.Body)
		return
	}
	defer r.Body.Close()

	// to prevent backend to timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// update the document that matches the buyerid and time of order
	_, err := db.EnergyBuyRequests.UpdateOne(
		ctx,
		bson.M{"reqid": reqId},
		bson.D{
			{"$set", bson.D{{"requestclosed", true}}},
		},
	)

	// if the update fails
	if err != nil {
		fmt.Println("Failed to close buy request")
		log.Fatal(err)
		return
	}
	respondWithJSON(w, r, http.StatusCreated, "Request Close Successful")
	return
}

/// for testing
func runLa() {
	///////////////////////////////
	// bruh
	currentBlockchain := getCurrentBlockchain()
	latestLocalBlock := currentBlockchain[12]
	fileDir := "/home/rafaquat/trn_sizes/trn_100.json"
	file, _ := json.MarshalIndent(latestLocalBlock, "", " ")

	_ = ioutil.WriteFile(fileDir, file, 0644)
	fmt.Println("Testing meeehh")

}
