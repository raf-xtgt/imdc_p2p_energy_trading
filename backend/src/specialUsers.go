package main

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"log"
	"net/http"
	"time"
)

/**All code concerning validators and clerks go here **/

func addValidator(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Methods", "GET,POST,OPTIONS,DELETE,PUT")

	var newValidator Validator
	var response SignUpResponse

	// get the data from json body
	decoder := json.NewDecoder(r.Body)
	// place the user data inside newUser
	if err := decoder.Decode(&newValidator); err != nil {
		fmt.Println("Failed adding a new user", err)
		respondWithJSON(w, r, http.StatusBadRequest, r.Body)
		return
	}
	defer r.Body.Close()

	// to prevent backend to timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// validate user
	if checkUsernameAndPass(newValidator.UserName, newValidator.Email) {

		// hash user password
		userPass := newValidator.Password
		newValidator.Password = hashPassword(userPass)

		// generate private and public key for user
		bitSize := 4096

		// Generate RSA private key.
		privateKey, err := rsa.GenerateKey(rand.Reader, bitSize)
		if err != nil {
			fmt.Println("Yo: Error generating keys", err.Error)
			panic(err)
		}

		// add the public key for the user
		pub := privateKey.Public()
		pubPEM := pem.EncodeToMemory(
			&pem.Block{
				Type:  "RSA PUBLIC KEY",
				Bytes: x509.MarshalPKCS1PublicKey(pub.(*rsa.PublicKey)),
			},
		)
		newValidator.PublicKey = pubPEM

		// Encode private key to PKCS#1 ASN.1 PEM.
		privatePEM := pem.EncodeToMemory(
			&pem.Block{
				Type:  "RSA PRIVATE KEY",
				Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
			},
		)
		newValidator.PrivateKey = privatePEM
		//write user info to the users collection
		writeUser, err := db.Users.InsertOne(ctx, newValidator)
		if err != nil {
			log.Fatal(err)
		}

		if addUniqueUserId(newValidator.UserName, newValidator.Email, fmt.Sprintf("%v", writeUser.InsertedID)) {
			//response.Res = true
			fmt.Println("New validator added after updating for id", writeUser.InsertedID)
			response.Res = true
			response.Email = newValidator.Email
			response.Validator = newValidator
			respondWithJSON(w, r, http.StatusCreated, response)
			return
			//fmt.Println("Validator Data from frontend", newValidator)
		}

	} else {
		fmt.Println("New user not added")
		response.Res = false
		respondWithJSON(w, r, http.StatusCreated, response)
		return

	}

}
