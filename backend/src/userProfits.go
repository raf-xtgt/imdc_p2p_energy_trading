package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func getUserIncomeData(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Methods", "GET,POST,OPTIONS,DELETE,PUT")
	var userId string
	var response Income

	// get the data from the request body
	decoder := json.NewDecoder(r.Body)
	// place the user data inside newUser
	if err := decoder.Decode(&userId); err != nil {
		fmt.Println("Failed adding new forecast data", err)
		respondWithJSON(w, r, http.StatusBadRequest, r.Body)
		return
	}
	defer r.Body.Close()
	fmt.Println("This is the user id for getting income data", userId)
	// to prevent backend to timeout
	mongoparams.ctx, mongoparams.cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer mongoparams.cancel()

	// get the current blockchain
	blockchain := getCurrentBlockchain()
	//transaction := getTransaction("626427e8ffa53e3720da2dd6")
	//fmt.Println(transaction)
	var receivables []float64
	var soldEnergy []float64
	var dates []string
	var tIds []string
	for i := 1; i < len(blockchain); i++ {
		// transactions in the block
		blockTrns := blockchain[i].Data
		hash := blockchain[i].Hash
		totalReceivable := 0.0
		totalEnSales := 0.0

		// loop through each transaction in the block
		for j := 0; j < len(blockTrns); j++ {
			// get the id of the transaction
			trnId := blockTrns[j].TId

			// get this transaction
			transaction := getTransaction(trnId)
			trnDate := transaction.Date

			trnBids := transaction.AuctionBids
			//fmt.Println(trnBids)
			//break
			//loop through the bids in this transaction
			for k := 0; k < len(trnBids); k++ {
				bidInfo := trnBids[k]
				bidSeller := bidInfo.SellerId

				// if this is the seller who sent the request
				if bidSeller == userId {
					fmt.Println("Bids that seller is in", bidInfo)
					totalReceivable += bidInfo.OptSellerReceivable
					totalEnSales += bidInfo.OptEnFromSeller
					dates = append(dates, trnDate)
					//receivables = append(receivables, bidInfo.OptSellerReceivable)
					//soldEnergy = append(soldEnergy, bidInfo.OptEnFromSeller)
					// tIds = append(tIds, hash)

				}

			}

		}
		// after looping through all the transactions for this block we append the data to the receivables
		receivables = append(receivables, totalReceivable)
		soldEnergy = append(soldEnergy, totalEnSales)
		tIds = append(tIds, hash)

	}
	// after looping through all transactions in the block where the seller was potentially present
	//var income Income
	response.Receivable = receivables
	response.EnergySold = soldEnergy
	response.BlockHashes = tIds
	response.Dates = dates
	//response = append(response, income)

	respondWithJSON(w, r, http.StatusCreated, response)
}
