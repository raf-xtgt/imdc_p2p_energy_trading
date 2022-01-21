package main

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"log"
	"net/http"
	"time"

	"crypto/sha256" //crypto library to hash the data
	"crypto/x509"

	// to create the secret(private) key for jwt token
	"crypto/rand"
	"crypto/rsa"
	"os"

	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson"
)

var secretkey string = ""

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
	var response LoginResponse
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Methods", "GET,POST,OPTIONS,DELETE,PUT")

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
		validToken, err := generateJWT(Profiles[0], "user")
		if err != nil {

			fmt.Println(err, "Failed to generate token")
			// w.Header().Set("Content-Type", "application/json")
			// json.NewEncoder(w).Encode(err)
			respondWithJSON(w, r, http.StatusBadRequest, r.Body)
			return
		}
		// create the jwt token
		var token Token
		//token.Email = Profiles[0].Email
		//token.Role = "user"
		token.TokenString = validToken
		response.Token = token
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(token)
		//respondWithJSON(w, r, http.StatusCreated, response)
	} else {
		fmt.Println("Username and password do not match")
	}
	return
}

/*
The GenerateJWT() function takes email and role as input.
Creates a token by HS256 signing method and adds authorized email, role,
and expiration time into claims. Claims are pieces of information added into tokens.
*/
func generateJWT(user User, role string) (string, error) {
	//secretkey := generateSecretKey()
	// secret key is the hash of the user's password
	secretkey = user.Password
	var mySigningKey = []byte(secretkey)
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	claims["authorized"] = true
	//claims["email"] = email
	claims["role"] = role
	claims["userInfo"] = user
	claims["exp"] = time.Now().Add(time.Minute * 30).Unix()

	tokenString, err := token.SignedString(mySigningKey)

	if err != nil {
		fmt.Errorf("Something Went Wrong: %s", err.Error())
		return "", err
	}
	return tokenString, nil
}

// Generate the secret key for jwt authorization, the key is an array of bytes
func generateSecretKey() []byte {
	// generate key pairs
	secretkey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		fmt.Printf("Cannot generate RSA key\n")
		os.Exit(1)
	}
	//fmt.Println("The secret key is", secretkey)

	// dump private key to file (store it in the os)
	var privateKeyBytes []byte = x509.MarshalPKCS1PrivateKey(secretkey)
	privateKeyBlock := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privateKeyBytes,
	}
	privatePem, err := os.Create("private.pem")
	if err != nil {
		fmt.Printf("error when create private.pem: %s \n", err)
		os.Exit(1)
	}
	err = pem.Encode(privatePem, privateKeyBlock)
	if err != nil {
		fmt.Printf("error when encode private pem: %s \n", err)
		os.Exit(1)
	}
	return privateKeyBytes
}

// function to check the integrity of the jwt that is sent from server
func isAuthorized(w http.ResponseWriter, r *http.Request) {

	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Methods", "GET,POST,OPTIONS,DELETE,PUT")

	var userToken Token
	fmt.Println(r.Header)
	if r.Header["Token"] == nil {
		fmt.Println("No Token Found")
		return
	}

	// get the data from json body
	decoder := json.NewDecoder(r.Body)
	// place the user data inside newUser
	if err := decoder.Decode(&userToken); err != nil {
		fmt.Println("Failed decoding token", err)
		respondWithJSON(w, r, http.StatusBadRequest, r.Body)
		return
	}
	defer r.Body.Close()
	var mySigningKey = []byte(secretkey)

	token, err := jwt.Parse(r.Header["Token"][0], func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("There was an error in parsing")
		}
		return mySigningKey, nil
	})

	if err != nil {
		fmt.Println("Your Token has been expired")
		json.NewEncoder(w).Encode(err)
		return
	}

	// check token validity
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if claims["role"] == "admin" {

			r.Header.Set("Role", "admin")
			return

		} else if claims["role"] == "user" {
			fmt.Println("User has been verified")
			r.Header.Set("Role", "user")
			respondWithJSON(w, r, http.StatusCreated, userToken)
			return
		}
	}
	return
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

/*
- In order to verify the jwt we need the secret key.
	- We need to store the secret key in a place where it can be retrieved at any time
	- privateKeyBytes == is what we need.
	- instead use the hash as the secret key lol
- Once we verify the jwt, we collect data from the claims part.
- The claims part will have the doc.id of the user in mongo db.
- We use the doc.id to get the required user data.
-

*/
