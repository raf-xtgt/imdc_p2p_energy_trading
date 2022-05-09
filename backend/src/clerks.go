package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

/**
This file contains:
- get all normal users who can be made into clerks
- convert a normal user into a clerk
- clerk validation
**/

const TOTAL_CLERKS = 2     // as of now there are only three clerks in the whole network
const BLOCKS_FOR_CHECK = 5 // need 5 new blocks for the integrity check

// function to get all normal users who can be made into clerks
func getAllUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Methods", "GET,POST,OPTIONS,DELETE,PUT")

	// to prevent backend to timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// filter all normal users
	cursor, err := db.Users.Find(ctx, bson.M{"type": "normal"})
	if err != nil {
		fmt.Println("Failed to get all users, clerks.go ::26")
		log.Fatal(err)
	}
	var profiles []User
	if err = cursor.All(ctx, &profiles); err != nil {
		log.Fatal(err)
	}
	fmt.Println("List of all users", len(profiles))
	respondWithJSON(w, r, http.StatusCreated, profiles)

}

// convert a normal user into a clerk
func convertToClerk(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Methods", "GET,POST,OPTIONS,DELETE,PUT")

	// to prevent backend to timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var userId string
	// get the data from json body
	decoder := json.NewDecoder(r.Body)
	// place the user data inside newUser
	if err := decoder.Decode(&userId); err != nil {
		fmt.Println("Failed getting userId in clerks.go ::56", err)
		respondWithJSON(w, r, http.StatusBadRequest, r.Body)
		return
	}
	defer r.Body.Close()

	//fmt.Println("UserId as received from frontend", userId)
	// to prevent backend to timeout

	// update the type to clerk that matches the user id
	_, err := db.Users.UpdateOne(
		ctx,
		bson.M{"uId": userId},
		bson.D{
			{"$set", bson.D{{"type", "clerk"}}},
		},
	)

	// if the update fails
	if err != nil {
		fmt.Println("Failed  updating userId in clerks.go ::77", err)
		log.Fatal(err)

	}
	fmt.Println("Made clerk successfully")
}

// perform the integrity check
func integrityCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Methods", "GET,POST,OPTIONS,DELETE,PUT")

	var response IntegrityCheckResponse
	// if five new blocks are present
	if newFiveBlocksExist() {
		// get the updated blockchain and user accounts for local copies
		createLocalCopies()

		// get the latest block from the central database
		latestBlock := getLatestBlock()

		// get the transactions in the latest block
		blockTransactions := latestBlock.Data

		// get the block info for the latest block --> hash, list of validators who checked it, blockId and list of clerks who checked it
		latestBlockMetadata := getBlockMetadata(latestBlock.Hash)

		// number of clerks who checked the block
		numClerkChecks := len(latestBlockMetadata.Clerks) - 1 // -1 to account for empty spot

		var fiftyPercent int = 50

		// if less than 50% of clerks have checked the latest block
		if (numClerkChecks/TOTAL_CLERKS)*100 < fiftyPercent {
			var counter = 0
			for j := 0; j < len(blockTransactions); j++ {
				transaction := blockTransactions[j]

				// verify this transaction using the local copy
				if localTrnVerification(transaction) {
					// validator checks the block only when they finish verifying all the transactions
					if counter == len(blockTransactions)-1 {
						fmt.Println("All transactions in latest block are checked by clerk")
						// use the nonce of the latest block and check whether its hash matches or not
						if checkBlock(latestBlock) {
							// add the validator in the list of validators who checked the block in blockInfo collection.
							updateCheckedClerks(latestBlock.Hash)
							response.IntegrityBreached = false
							respondWithJSON(w, r, http.StatusCreated, response)
						} else {
							fmt.Println("Block hash doesn't match, then need to find culprit validator")
							response.IntegrityBreached = true
							respondWithJSON(w, r, http.StatusCreated, response)
						}

					} else {
						counter += 1
					}
				}
			}
		} else {
			fmt.Println(" More than 50 percent of clerks have checked the block, so chill")

			response.IntegrityBreached = false
			respondWithJSON(w, r, http.StatusCreated, response)
		}

	} else {
		fmt.Println("No need to do integrity check since there are no 5 new blocks")
		response.IntegrityBreached = false
		respondWithJSON(w, r, http.StatusCreated, response)
	}

}

// function to update the list of clerks who checked the block with no repeat clerks
func updateCheckedClerks(hash string) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := db.BlockInfo.Find(ctx,
		bson.M{"hash": hash})

	if err != nil {
		log.Fatal(err)
		fmt.Println("Failed to retrieve block info :localBlockchain.go 199")
	}

	var blockData []BlockInfo
	if err = cursor.All(ctx, &blockData); err != nil {
		log.Fatal(err)
		fmt.Println("Failed to write the block metadata :localBlockchain.go 205")
	}

	var found bool
	data := blockData[0]
	blockClerks := data.Clerks

	for i := 0; i < len(blockClerks); i++ {
		vId := blockClerks[i]
		//		fmt.Println("Checked validator", vId, "  Logged in user", loggedInUser)
		if loggedInUser == vId {
			found = true
		} else if i == len(blockClerks) {
			// checked all values and not found validator
			found = false
			// then we add the validator
			blockClerks = append(blockClerks, loggedInUser)
		}
	}

	// if the current validator is not part of the list, then update the list
	if !found {

		_, err := db.BlockInfo.UpdateOne(
			ctx,
			bson.M{"hash": hash},
			bson.D{
				{"$push", bson.D{{"clerks", loggedInUser}}},
			},
		)

		// if the update fails
		if err != nil {
			log.Fatal(err)
			fmt.Println("Failed to update the block metadata on successful block check")
		} else {
			fmt.Println("Successfully updated the block metadata")
		}

	} else {
		fmt.Println("Current user checked the latest block ady")
	}

}

//function to check for 5 new blocks
func newFiveBlocksExist() bool {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// get the latest index. The collection always has only one document
	cursor, err := db.LatestIndex.Find(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	var latestIndex []LatestIndex
	if err = cursor.All(ctx, &latestIndex); err != nil {
		log.Fatal(err)
	}

	if (latestIndex[0].NewBlockNum % BLOCKS_FOR_CHECK) == 0 {
		return true
	} else {
		return false
	}
}
