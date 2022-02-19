package main

import (
	"context"
	"fmt"
	"log"
	"time"
)

/**
Function to create a database entry for user balance
ASSUMPTION 1: The energy balance comes from the smart metre
ASSUMPTION 2: We assume already in sync with some financial institution
**/
func createUserAccount(userId string) {
	// to prevent backend to timeout
	mongoparams.ctx, mongoparams.cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer mongoparams.cancel()
	var balance AccountBalance
	balance.UserId = userId
	// add dummy data
	balance.FiatBalance = 2000
	balance.EnergyBalance = 80

	//write the balance in the database info to the users collection
	writeData, err := db.UserAccBalance.InsertOne(mongoparams.ctx, balance)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("User balance added", writeData.InsertedID)
	return

}
