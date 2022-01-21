package main

import (
	"context"
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

// MonogParam instance to make the context, client and cancel to be available gloablly
var mongoparams MongoParam

// MonogDatabase instance to make the cluster and collections to be available gloablly
var db MongoDatabase

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
	mongoparams.ctx = ctx
	mongoparams.cancel = cancel
	mongoparams.client = client
	connectToDb()
	log.Fatal(listen())
}

func connectToDb() MongoDatabase {

	// use the db params to access the cluster
	db.Cluster = mongoparams.client.Database("IMDC-p2p-energy")
	// refer to the users collection
	db.Users = db.Cluster.Collection("users")
	return db
}

func getEnvVar(key string) string {
	// load .env file
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}

func listen() error {
	//http.HandleFunc("/")
	// when a request is made on/register, then run addNewUser function
	http.HandleFunc("/Register", addNewUser)
	http.HandleFunc("/GetUser", getUser)
	http.HandleFunc("/Login", authenticateUser)
	http.HandleFunc("/VerifyToken", isAuthorized)
	log.Fatal(http.ListenAndServe(":8080", nil))

	return nil
}

func getUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Methods", "GET,POST,OPTIONS,DELETE,PUT")
}
