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

type MongoDatabases struct {
	Blockchain *mongo.Database
	Users      *mongo.Database
	Market     *mongo.Database
}

// instance of MongoParam
var mongoparameters MongoParam
var MongoDBs MongoDatabases

func main() {
	fmt.Println("Starting server")
	// get the environment variables to required for database authentication
	dbCluster := getEnvVar("DB_CLUSTER_ADDR")
	dbUserName := getEnvVar("DB_USERNAME")
	dbPass := getEnvVar("DB_PASSWORD")
	fmt.Println(dbCluster)

	// Establish connection to mongodb cluster
	dbUrl := "mongodb+srv://" + dbUserName + ":" + dbPass + "@imdc-p2p-energy.y0a68.mongodb.net/" + dbCluster + "?retryWrites=true&w=majority"

	clientOptions := options.Client().
		ApplyURI(dbUrl)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to backend server")

	// set the db parameters to the global variables
	mongoparameters.ctx = ctx
	mongoparameters.cancel = cancel
	mongoparameters.client = client
	connectToDb("all")
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

func connectToDb(Choice string) *mongo.Database {

	BlockchainDatabase := mongoparameters.client.Database("IMDC-p2p-energy")
	UserDatabase := mongoparameters.client.Database("users")

	var Database *mongo.Database
	if Choice == "Blockchain" {
		Database = BlockchainDatabase
	} else if Choice == "Users" {
		Database = UserDatabase
	}
	MongoDBs.Blockchain = BlockchainDatabase
	MongoDBs.Users = UserDatabase
	return Database
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
	if err := decoder.Decode(&NewUser); err != nil {
		fmt.Println("Error")
		fmt.Println(err)
		respondWithJSON(w, r, http.StatusBadRequest, r.Body)
		return
	}
	defer r.Body.Close()
	fmt.Println("This is the data from frontend", NewUser)

	// // // handle to the database cluster
	// database := mongoparameters.client.Database("IMDC-p2p-energy")
	// // // handle to the collection inside the cluster
	// usersCollection := database.Collection("users")
	// // // write data to the users collection
	// writeUser, err := usersCollection.InsertOne(mongoparameters.ctx, NewUser)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println("Id of user2", writeUser.InsertedID)
	// respondWithJSON(w, r, http.StatusCreated, NewUser)
	// return
}

func listen() error {
	//http.HandleFunc("/")
	// when a request is made on/register, then run addNewUser function

	http.HandleFunc("/", addNewUser)
	log.Fatal(http.ListenAndServe(":8080", nil))

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
