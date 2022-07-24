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
	PrivateKey   []byte `bson:privateKey`
	Address      string `bson:"address"`
	UId          string `bson:id`
	SmartMeterNo int    `bson:"smartMeterNo"`
	Type         string `bson:type` // normal, validator, clerk
}

// represents validator and one administrator
type Validator struct {
	UserName   string `bson:"username"`
	Email      string `bson:"email"`
	Password   string `bson:"password"`
	PublicKey  []byte `bson:publicKey`
	PrivateKey []byte `bson:privateKey`
	Address    string `bson:"address"`
	UId        string `bson:id`
	Type       string `bson:type` // normal, validator, clerk
	FullName   string `bson:"fullName"`
	ICNum      int    `bson:identificationCardNo`
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
	Blockchain         *mongo.Collection
	Transactions       *mongo.Collection
	LatestIndex        *mongo.Collection //collection to hold the latest index
	Trigger            *mongo.Collection
	BlockInfo          *mongo.Collection // collection to hold the block Id and the number of validators who checked it
}

type BlockInfo struct {
	BlockId    string   "json:blockId"
	Validators []string "json:validators" // the list of validators who have checked the block
	Hash       string   "json:hash"       //hash of this block
	Clerks     []string "json:clerks"     // list of clerks who have checked the block. Clerks will check after every five new blocks are added.
}

// Structure that will be sent as sign up response to frontend
type SignUpResponse struct {
	Email     string
	Res       bool
	User      User
	Validator Validator
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

// structure of a single block
type Block struct {
	Index    int           "json: index"
	Data     []Transaction "json: data"
	Hash     string        "json: hash"
	PrevHash string        "json: prevHash"
	Nonce    string        "json: nonce"
}

// structure of a single transaction
type Transaction struct {
	BuyerId                      string  "json: buyerId"
	BuyerPayable                 float64 "json: buyerPayable"
	BuyerEnReceivableFromAuction float64 "json: buyerEnReceivableFromAuction"
	BuyerEnReceivableFromTNB     float64 "json: buyerEnReceivableFromTNB"
	AuctionBids                  []Bid   "json: auctionBids"
	TNBReceivable                float64 "json: TNBReceivable"
	Verified                     bool    "json:verified"
	TId                          string  "json: tId"
	Checks                       int     "json: checks"
	Date                         string  "json: date"
}

// structure of a bid made by seller in the auction
type Bid struct {
	SellerId            string  "json: sellerId"
	OptEnFromSeller     float64 "json: optEnFromSeller"
	OptSellerReceivable float64 "json: optSellerReceivable"
	SellerFiatBalance   float64 "json: sellerFiatBalance"
	SellerEnergyBalance float64 "json: sellerEnergyBalance"
}

type Trigger struct {
	NewBlockExists bool "json: newBlockExists"
}

type BlockchainResponse struct {
	Blockchain []Block
}

// response sent to frontedn
type IntegrityCheckResponse struct {
	IntegrityBreached bool "json: integrityBreached"
}

// response to send user income data
type Income struct {
	Receivable  []float64 "json: receivable" // amount of money the seller is going to receive
	EnergySold  []float64 "json: energySold" // amount of energy sold by this seller for the receivable
	Payable     []float64 "json:payable"     // amount of money this user paid for the energy
	BlockHashes []string  "json: blockHash"  // id of transaction on which the bid was made
	Dates       []string  "json: dates"
}

type TransactionsResponse struct {
	Transactions []Transaction "json: transactions"
}
