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
	Address      string `bson:"address"`
	UId          string `bson:id`
	SmartMeterNo int    `bson:"smartMeterNo"`
}

// userId structure that represents a unique user id
// type UserId struct {
// 	Id bson.ObjectId `bson:"_id,omitempty"` // only uppercase variables can be exported
// }

type MongoDatabase struct {
	Cluster           *mongo.Database   //cluster
	Users             *mongo.Collection //collection
	EnergyPriceHouse  *mongo.Collection
	EnergyBuyRequests *mongo.Collection
	BuyOrderForecast  *mongo.Collection
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
}

// structure of response when buy energy requests are queried
type GetBuyEnergyResponse struct {
	Requests []BuyRequest "json:buyEnergyRequests"
}

// type BuyEnergyForecast struct {
// 	ActualX    []string  "json:actual_x" // x axis data
// 	ActualY    []float64 "json:actual_y"
// 	Pred_X     []string  "json:pred_x"
// 	Pred_Y     []float64 "json:pred_y"
// 	Prediction float64   "json:current_pred" //the prediction for the next 30 minutes of time(future 30 min)
// 	DateStr    string    "json:date"
// }

type BuyForecastResponse struct {
	Actual_X     []string  "json:actual_x" // x axis data
	Actual_Y     []float64 "json:actual_y"
	Pred_X       []string  "json:pred_x"
	Pred_Y       []float64 "json:pred_y"
	Current_Pred float64   "json:current_pred" //the prediction for the next 30 minutes of time(future 30 min)
	Date         string    "json:date"
}
