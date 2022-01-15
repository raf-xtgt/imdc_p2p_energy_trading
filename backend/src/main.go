package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	// Import godotenv
	"github.com/joho/godotenv"

	// To connect to mongodb
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// globally define the db parameters so that other methods can access the db
type MongoParam struct {
	ctx    context.Context
	cancel context.CancelFunc
	client *mongo.Client
}

type User struct {
	FullName string `bson:"name"`
	Email    string `bson:"email"`
}

// instance of MongoParam
var mongoparameters MongoParam

func main() {
	fmt.Println("Starting server")
	// get the environment variables to required for database authentication
	dbCluster := getEnvVar("DB_CLUSTER_ADDR")
	dbUserName := getEnvVar("DB_USERNAME")
	dbPass := getEnvVar("DB_PASSWORD")
	fmt.Println(dbCluster)

	// Establish connection to mongodb cluster
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb+srv://" + dbUserName + ":" + dbPass + "@imdc-p2p-energy.y0a68.mongodb.net/myFirstDatabase?retryWrites=true&w=majority"))
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)

	// ping the cluster
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to backend server")

	// set the db parameters to the global variables
	mongoparameters.ctx = ctx
	mongoparameters.cancel = cancel
	mongoparameters.client = client
	log.Fatal(listen())
	// // handle to the database cluster
	// database := client.Database("IMDC-p2p-energy")
	// // handle to the collection inside the cluster
	// usersCollection := database.Collection("users")
	// // write data to the users collection
	// writeUser, err := usersCollection.InsertOne(ctx, bson.D{
	// 	{"title", "rafaquat"},
	// 	{"arrays", bson.A{"This", "is", "an", "array"}},
	// })
	// fmt.Println("Id of the user", writeUser.InsertedID)

}

func getEnvVar(key string) string {
	// load .env file
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}

func addNewUser(w http.ResponseWriter, r *http.Request) {
	var NewUser User
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Methods", "GET,POST,OPTIONS,DELETE,PUT")

	// get the data from json body
	decoder := json.NewDecoder(r.Body)
	fmt.Println(r.Body)
	if err := decoder.Decode(&NewUser); err != nil {
		fmt.Println("Got error lol")
		fmt.Println(err)
		respondWithJSON(w, r, http.StatusBadRequest, r.Body)
		return
	}
	defer r.Body.Close()
	fmt.Println("This is the data from frontend", NewUser)

	// // handle to the database cluster
	database := mongoparameters.client.Database("IMDC-p2p-energy")
	// // handle to the collection inside the cluster
	usersCollection := database.Collection("users")
	// // write data to the users collection
	writeUser, err := usersCollection.InsertOne(mongoparameters.ctx, NewUser)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Id of user2", writeUser.InsertedID)
	respondWithJSON(w, r, http.StatusCreated, NewUser)
}

func listen() error {
	//http.HandleFunc("/")
	// when a request is made on/register, then run addNewUser function
	http.HandleFunc("/register", addNewUser)

	log.Fatal(http.ListenAndServe(":"+getEnvVar("PORT"), nil))
	return nil
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
