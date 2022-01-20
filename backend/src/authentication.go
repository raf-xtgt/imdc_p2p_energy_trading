package main

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"crypto/sha256" //crypto library to hash the data

	"go.mongodb.org/mongo-driver/bson"
)

func addNewUser(w http.ResponseWriter, r *http.Request) {
	var newUser User
	var response AuthResponse
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Methods", "GET,POST,OPTIONS,DELETE,PUT")

	// get the data from json body
	decoder := json.NewDecoder(r.Body)
	// place the user data inside newUser
	if err := decoder.Decode(&newUser); err != nil {
		fmt.Println("Failed adding a new user", err)
		respondWithJSON(w, r, http.StatusBadRequest, r.Body)
		return
	}
	defer r.Body.Close()
	//fmt.Println("This is the data from frontend", NewUser)
	// to prevent backend to timeout
	mongoparams.ctx, mongoparams.cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer mongoparams.cancel()

	// validate user
	if checkUsernameAndPass(newUser.UserName, newUser.Email) {
		// hash user password
		userPass := newUser.Password
		hash := sha256.New()
		hash.Write([]byte(userPass))
		passHash := hash.Sum(nil)
		newUser.Password = hex.EncodeToString(passHash)

		//write user info to the users collection
		writeUser, err := db.Users.InsertOne(mongoparams.ctx, newUser)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("New user added", writeUser.InsertedID)
		response.Res = true

	} else {
		fmt.Println("New user not added")
		response.Res = false
	}

	response.User = newUser
	response.Email = newUser.Email
	respondWithJSON(w, r, http.StatusCreated, response)
	//respondWithJSON(w, r, http.StatusCreated, NewUser)
	return
}

// function to check whether the email and username already exist or not
func checkUsernameAndPass(username string, email string) (result bool) {
	result = false
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// find username
	cursor, err := db.Users.Find(ctx, bson.M{"username": username, "email": email})
	if err != nil {
		log.Fatal(err)
	}

	var Profiles []User
	if err = cursor.All(ctx, &Profiles); err != nil {
		log.Fatal(err)

	}
	if len(Profiles) == 0 {
		fmt.Println("No users with this username")
		result = true
	} else {
		fmt.Println("Username and email already being used")
		result = false
	}
	fmt.Println(Profiles)
	return result

	// defer cursor.Close(ctx)
	// for cursor.Next(ctx) {
	// 	var episode bson.M
	// 	if err = cursor.Decode(&episode); err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	fmt.Println("Data:", episode)
	// }
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
