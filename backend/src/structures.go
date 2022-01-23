package main

import (
	"context"

	jwt "github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/mongo"
)

// globally define the db parameters so that other methods can access the db
type MongoParam struct {
	ctx    context.Context
	cancel context.CancelFunc
	client *mongo.Client
}

type User struct {
	UserName     string `bson:"username"`
	Email        string `bson:"email"`
	Password     string `bson:"password"`
	Address      string `bson:"address"`
	SmartMeterNo int    `bson:"smartMeterNo"`
}

type MongoDatabase struct {
	Cluster *mongo.Database   //cluster
	Users   *mongo.Collection //collection
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
	// https://tools.ietf.org/html/rfc7519
	jwt.StandardClaims
	CustomClaims map[string]string `json:"custom,omitempty"`
}
