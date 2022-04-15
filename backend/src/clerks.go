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
This file contains:
- clerk validation
- get all normal users who can be made into clerks
**/

// function to get all normal users who can be made into clerks
func getAllUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Methods", "GET,POST,OPTIONS,DELETE,PUT")

	// to prevent backend to timeout
	mongoparams.ctx, mongoparams.cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer mongoparams.cancel()

	// var response []User

	cursor, err := db.Users.Find(mongoparams.ctx, bson.M{"type": "normal"})
	if err != nil {
		fmt.Println("Failed to get all users, clerks.go ::26")
		log.Fatal(err)
	}
	var profiles []User
	if err = cursor.All(mongoparams.ctx, &profiles); err != nil {
		log.Fatal(err)
	}
	fmt.Println("List of all users", len(profiles))
	respondWithJSON(w, r, http.StatusCreated, profiles)

}

func convertToClerk(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Methods", "GET,POST,OPTIONS,DELETE,PUT")

	// to prevent backend to timeout
	mongoparams.ctx, mongoparams.cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer mongoparams.cancel()

	var userId string
	// get the data from json body
	decoder := json.NewDecoder(r.Body)
	// place the user data inside newUser
	if err := decoder.Decode(&userId); err != nil {
		fmt.Println("Failed getting userId in clerks.go ::56", err)
		respondWithJSON(w, r, http.StatusBadRequest, r.Body)
		return
	}
	defer r.Body.Close()

	//fmt.Println("UserId as received from frontend", userId)
	// to prevent backend to timeout
	mongoparams.ctx, mongoparams.cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer mongoparams.cancel()

	// update the type to clerk that matches the user id
	_, err := db.Users.UpdateOne(
		mongoparams.ctx,
		bson.M{"uId": userId},
		bson.D{
			{"$set", bson.D{{"type", "clerk"}}},
		},
	)

	// if the update fails
	if err != nil {
		fmt.Println("Failed  updating userId in clerks.go ::77", err)
		log.Fatal(err)

	}
	fmt.Println("Made clerk successfully")
}
