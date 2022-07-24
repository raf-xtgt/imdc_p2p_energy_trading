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

	// for
	"github.com/rs/cors"
)

// MonogParam instance to make the context, client and cancel to be available gloablly
var mongoparams MongoParam

// MonogDatabase instance to make the cluster and collections to be available gloablly
var db MongoDatabase

func main() {
	//generateEnergyPriceData()
	//processTime()
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
	db.EnergyPriceHouse = db.Cluster.Collection("householdEnergyPrice")
	db.EnergyBuyRequests = db.Cluster.Collection("buyRequests")
	db.EnergySellRequests = db.Cluster.Collection("sellRequests")
	db.BuyOrderForecast = db.Cluster.Collection("buyOrderForecast")
	db.ProdForecast = db.Cluster.Collection("energy_forecast")
	db.UserAccBalance = db.Cluster.Collection("accountBalance")
	db.Blockchain = db.Cluster.Collection("blockchain")
	db.Transactions = db.Cluster.Collection("transactions")
	db.LatestIndex = db.Cluster.Collection("latestIndex")
	db.Trigger = db.Cluster.Collection("trigger") // document holds boolean value when a new block is made
	db.BlockInfo = db.Cluster.Collection("blockInfo")
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
	mux := http.NewServeMux()
	mux.HandleFunc("/AddHouseholdData", createEnergyData)
	// when a request is made on/register, then run addNewUser function
	mux.HandleFunc("/Register", addNewUser)
	mux.HandleFunc("/GetUserDetails", getUserDetails)
	mux.HandleFunc("/Login", authenticateUser)
	mux.HandleFunc("/VerifyToken", isAuthorized)
	mux.HandleFunc("/CreateBuyRequest", addBuyRequest)
	mux.HandleFunc("/CreateSellRequest", addSellRequest)
	mux.HandleFunc("/GetBuyRequests", getEnergyRequests)
	mux.HandleFunc("/RunBuyEnergyForecast", runBuyEnergyForecast)
	mux.HandleFunc("/GetLatestBuyForecast", getLatestBuyForecast)
	mux.HandleFunc("/RunSellEnergyForecast", runSellEnergyForecast)
	mux.HandleFunc("/GetLatestSellForecast", getLatestSellForecast)
	mux.HandleFunc("/CloseBuyRequest", closeBuyRequest)
	mux.HandleFunc("/RunDoubleAuction", runDoubleAuction)
	mux.HandleFunc("/AddValidator", addValidator)
	mux.HandleFunc("/CreateGenesisBlock", createGenesisBlock)
	mux.HandleFunc("/UpdateBlockchain", updateChain)
	mux.HandleFunc("/GetBlockchain", sendBlockchainToFrontend)
	mux.HandleFunc("/GetAllUsers", getAllUsers)
	mux.HandleFunc("/MakeClerk", convertToClerk)
	mux.HandleFunc("/ClerkIntegrityCheck", integrityCheck)
	mux.HandleFunc("/GetUserIncome", getUserIncomeData)
	mux.HandleFunc("/GetTNBIncome", getTNBIncomeData)
	mux.HandleFunc("/GetUserBuyRequests", getUserBuyRequests)
	mux.HandleFunc("/GetUserSellRequests", getUserSellRequests)

	handler := cors.Default().Handler(mux)
	log.Fatal(http.ListenAndServe(":8080", handler))

	return nil
}

func getUserDetails(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Methods", "GET,POST,OPTIONS,DELETE,PUT")
	var request string

	// get the data from json body
	decoder := json.NewDecoder(r.Body)
	// place the user data inside newRequest
	if err := decoder.Decode(&request); err != nil {
		fmt.Println("Failed getting user data", err)
		respondWithJSON(w, r, http.StatusBadRequest, r.Body)
		return
	}
	defer r.Body.Close()

	// to prevent backend to timeout
	mongoparams.ctx, mongoparams.cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer mongoparams.cancel()

	userData, err := getAccountData(request)
	if err != nil {
		log.Println(err)
		http.Error(w, "Request failed!", http.StatusUnauthorized)
	}
	//fmt.Println("The user data as per id", userData)
	w.Write(userData)
}
