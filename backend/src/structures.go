package main

import (
	"context"

	jwt "github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/mongo"
	//"labix.org/v2/mgo/bson"
)

// globally define the db parameters so that other methods can access the db
type MongoParam struct {
	ctx    context.Context
	cancel context.CancelFunc
	client *mongo.Client
}

// structure to represent a user
type User struct {
	UserName     string `bson:"username"`
	Email        string `bson:"email"`
	Password     string `bson:"password"`
	PublicKey    []byte `bson:publicKey`
	Address      string `bson:"address"`
	UId          string `bson:id`
	SmartMeterNo int    `bson:"smartMeterNo"`
}

// structure to represent a user account balance
type AccountBalance struct {
	UserId        string  "bson: userId"
	FiatBalance   float64 "bson: fiatBalance"
	EnergyBalance float64 "bson: energyBalance"
}

// lists all the collections in the db
type MongoDatabase struct {
	Cluster            *mongo.Database   //cluster
	Users              *mongo.Collection //collection
	EnergyPriceHouse   *mongo.Collection
	EnergyBuyRequests  *mongo.Collection
	EnergySellRequests *mongo.Collection
	BuyOrderForecast   *mongo.Collection
	ProdForecast       *mongo.Collection
	UserAccBalance     *mongo.Collection
}

// Structure that will be sent as sign up response to frontend
type SignUpResponse struct {
	Email string
	Res   bool
	User  User
}

// JWT structure upon login
type Token struct {
	//Role        string `json:"role"`
	//Email       string `json:"email"`
	TokenString string `json:"token"`
}

// Structure to be sent as login response to frontend
type LoginResponse struct {
	Token Token
}

type JWTData struct {
	// Standard claims are the standard jwt claims from the IETF standard
	jwt.StandardClaims
	CustomClaims map[string]string `json:"custom,omitempty"`
}

type JWTVerifiedData struct {
	//Email    string `json:"email"`
	//Username string `bson:"username"`
	User User
}

// structure to hold all the required data
type EnergyPriceData struct {
	Day      string      "json: day" // day of the week
	Average  float64     "json: average"
	Data     [24]float64 "json: data"
	DateStr  string      "json:dateString" // date in dd-mm-yyyy format
	DateTime int64       "json: dateTime"  // time in nano for sorting
}

// structure to respond with data to the frontend
type EnergyDataResponse struct {
	Data EnergyPriceData
}

// structure representing buy request coming from frontend
type BuyRequest struct {
	BuyerId       string  "json: buyerId"
	EnergyAmount  float64 "json: energyAmount"
	FiatAmount    float64 "json: fiatAmount"
	RequestClosed bool    "json: requestClosed"
	ReqTime       string  "json: requestTime"
	ReqId         string  "json:requestId"
	Auctioned     bool    "json:auctioned" // to check whether the request hsa undergone the double auction
}

// structure representing sell(bid) request coming from frontend
type SellRequest struct {
	SellerId     string  "json: sellerId"
	EnergyAmount float64 "json: energyAmount"
	FiatAmount   float64 "json: fiatAmount"
	SellReqId    string  "json:sellRequestId"
	ReqTime      string  "json: requestTime"
	BuyReqId     string  "json:buyReqId" //the id of the buy request on which the bid is made
}

// structure of response when buy energy requests are queried
type GetBuyEnergyResponse struct {
	Requests []BuyRequest "json:buyEnergyRequests"
}

// response when forecasting data is sent when user buys energy
type BuyForecastResponse struct {
	Actual_X     []string  "json:actual_x" // x axis data
	Actual_Y     []float64 "json:actual_y"
	Pred_X       []string  "json:pred_x"
	Pred_Y       []float64 "json:pred_y"
	Current_Pred float64   "json:current_pred" //the prediction for the next 30 minutes of time(future 30 min)
	Date         string    "json:date"
}

type ProdForecastRequest struct {
	UserId string "json: userId"
	Date   string "json: date"
}
