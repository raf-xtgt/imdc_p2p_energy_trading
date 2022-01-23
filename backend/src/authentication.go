package main

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"crypto/sha256" //crypto library to hash the data

	// to create the secret(private) key for jwt token

	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson"
)

//var secretkey string = ""
//var SECRET string = "42isTheAnswer"
var jwtSecret string = ""

func addNewUser(w http.ResponseWriter, r *http.Request) {
	var newUser User
	var response SignUpResponse
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
		newUser.Password = hashPassword(userPass)

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

	// find username and password
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
	//fmt.Println(Profiles)
	return result
}

/** Function to handle user login **/
func authenticateUser(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Authenticating User")
	var newUser User
	//var response LoginResponse
	//w.Header().Add("Access-Control-Allow-Origin", "*")
	//w.Header().Add("Access-Control-Allow-Methods", "GET,POST,OPTIONS,DELETE,PUT")

	// get the data from json body
	decoder := json.NewDecoder(r.Body)
	// place the user data inside newUser
	if err := decoder.Decode(&newUser); err != nil {
		fmt.Println("Failed adding a new user", err)
		// w.Header().Set("Content-Type", "application/json")
		// json.NewEncoder(w).Encode(err)
		respondWithJSON(w, r, http.StatusBadRequest, r.Body)
		return
	}
	defer r.Body.Close()
	//fmt.Println("This is the data from frontend", NewUser)
	// to prevent backend to timeout
	mongoparams.ctx, mongoparams.cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer mongoparams.cancel()

	// hash the password entered by user
	userPass := newUser.Password
	hashedPass := hashPassword(userPass) // if passwords match, then the hash should be the same

	cursor, err := db.Users.Find(mongoparams.ctx, bson.M{"username": newUser.UserName, "password": hashedPass})
	if err != nil {
		log.Fatal(err)
	}

	var Profiles []User
	if err = cursor.All(mongoparams.ctx, &Profiles); err != nil {
		log.Fatal(err)
	}

	// there should be only 1 profile with the given username and email
	if len(Profiles) == 1 {
		fmt.Println("Username and password match an account in the db")
		fmt.Println("Matched profile info", Profiles)
		jwtSecret = userPass

		// creating the JWT
		claims := JWTData{
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: time.Now().Add(time.Hour).Unix(),
			},
			CustomClaims: map[string]string{
				"username": Profiles[0].UserName,
				"email":    Profiles[0].Email,
			},
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString([]byte(jwtSecret))
		if err != nil {
			log.Println(err)
			http.Error(w, "Login failed!", http.StatusUnauthorized)
		}

		json, err := json.Marshal(struct {
			Token string `json:"token"`
		}{
			tokenString,
		})

		if err != nil {
			log.Println(err)
			http.Error(w, "Login failed!", http.StatusUnauthorized)
		}

		w.Write(json)
		return
	} else {
		fmt.Println("Username and password do not match")
	}
	return
}

// function to check the integrity of the jwt that is sent from server
func isAuthorized(w http.ResponseWriter, r *http.Request) {
	//w.Header().Add("Access-Control-Allow-Origin", "*")
	//w.Header().Add("Access-Control-Allow-Methods", "GET,POST,OPTIONS,DELETE,PUT")
	var jwtToken string

	// get the data from json body
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&jwtToken); err != nil {
		fmt.Println("Failed adding a new user", err)
		respondWithJSON(w, r, http.StatusBadRequest, r.Body)
		return
	}
	fmt.Println("Request body:\n", jwtToken)

	claims, err := jwt.ParseWithClaims(jwtToken, &JWTData{}, func(token *jwt.Token) (interface{}, error) {
		if jwt.SigningMethodHS256 != token.Method {
			return nil, errors.New("Invalid signing algorithm")
		}
		return []byte(jwtSecret), nil
	})

	if err != nil {
		log.Println(err)
		http.Error(w, "Request failed!", http.StatusUnauthorized)
	}

	data := claims.Claims.(*JWTData)

	userName := data.CustomClaims["username"]
	userEmail := data.CustomClaims["email"]
	jsonData, err := getAccountData(userName, userEmail)
	if err != nil {
		log.Println(err)
		http.Error(w, "Request failed!", http.StatusUnauthorized)
	}

	w.Write(jsonData)
}

func getAccountData(username string, email string) ([]byte, error) {
	mongoparams.ctx, mongoparams.cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer mongoparams.cancel()

	cursor, err := db.Users.Find(mongoparams.ctx, bson.M{"username": username, "email": email})
	if err != nil {
		log.Fatal(err)
	}

	var Profile []User
	if err = cursor.All(mongoparams.ctx, &Profile); err != nil {
		log.Fatal(err)
	}

	output := JWTVerifiedData{Profile[0].Email, Profile[0].UserName}
	json, err := json.Marshal(output)
	if err != nil {
		return nil, err
	}

	return json, nil
}

/** To return data as a json to the frontend */
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

// function to hash the password using sha256
func hashPassword(userPass string) string {
	hash := sha256.New()
	hash.Write([]byte(userPass))
	passHash := hash.Sum(nil)
	hashOutput := hex.EncodeToString(passHash)
	return hashOutput
}
