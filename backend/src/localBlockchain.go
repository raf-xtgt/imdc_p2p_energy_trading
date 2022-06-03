package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

/** This file has all the code to :
- store the blockchain locally
- update the locally stored chain
- read from the locally stored chain
**/

const LOCAL_BLOCKCHAIN = "blockchain.json"
const TOTAL_VALIDATORS = 2 // total number of validators in the whole network

// create a locally stored blockchain
func createLocalBlockchainFile(dirname string) string {

	filename := dirname + "/" + LOCAL_BLOCKCHAIN
	var _, err = os.Stat(filename)

	// create the local blockchain if it does not already exist
	if os.IsNotExist(err) {
		file, err := os.Create(filename)
		if err != nil {
			fmt.Println(err)

		}
		//fmt.Println("File created successfully", filename)
		defer file.Close()
	} else {
		//fmt.Println("File already exists!", filename)

	}

	return filename

}

// write the blockchain in the file
// we always write it so mo need to update. writing is updating itself
func writeLocalBlockchain(data []Block, fileDir string) {
	file, _ := json.MarshalIndent(data, "", " ")

	_ = ioutil.WriteFile(fileDir, file, 0644)
	fmt.Println("Local blockchain is now up to date")
}

// read the blocks from the local blockchain
func readLocalBlockchain(filepath string) []Block {
	file, _ := ioutil.ReadFile(filepath)

	localBlockchain := []Block{}

	_ = json.Unmarshal([]byte(file), &localBlockchain)

	//fmt.Println("Local blockchain as read", localBlockchain)
	return localBlockchain

}

// the validators who did not produce the block will need to veirfy transactions and sign the block
func verifyCentralBlockchain() bool {
	fmt.Println("inside verify central blockchain")
	//currentChain := getCurrentBlockchain()
	//latestBlock := getLatestBlock()
	//fmt.Println("length of current chain", len(currentChain))
	// var latestBlockArr [1]Block
	// latestBlockArr[0] = latestBlock
	// loop through blockchain but skip the genesis block
	latestCentralBlock := getLatestBlock()
	fmt.Println("\nBlock under check", latestCentralBlock)
	blockTransactions := latestCentralBlock.Data
	latestBlockMetadata := getBlockMetadata(latestCentralBlock.Hash)

	// we only check if the latest block was not checked by all validators
	// this means validators array in latestBlockMetadata would be less than the TOTAL_VALIDATORS
	if len(latestBlockMetadata.Validators) < TOTAL_VALIDATORS {
		var counter = 0
		for j := 0; j < len(blockTransactions); j++ {
			//fmt.Println("Inside block transactions")
			transaction := blockTransactions[j]

			// if the transaction is verified according to local copy
			if localTrnVerification(transaction) {
				// increment validator check in the central database
				incrementChecks(transaction.TId)

				// validator checks the block only when they finish verifying all the transactions
				if counter == len(blockTransactions)-1 {
					fmt.Println("All transactions in latest block are checked")
					// use the nonce of the latest block and check whether its hash matches or not
					if checkBlock(latestCentralBlock) {
						// add the validator in the list of validators who checked the block in blockInfo collection.
						updateCheckedValidators(latestCentralBlock.Hash)
						//return true
					}

				} else {
					counter += 1
				}

			} else {
				fmt.Println("Local transaction verification failed")
			}

		}
	} else {
		// no need to check since both validators have done so
		fmt.Println("no need to check since both validators have done ")
		//setTrigger(false) // no new blocks to check so make it false
		//return true
	}

	return true

}

// other not currently logged in validators check whether user has sufficient balance or not
func localTrnVerification(transaction Transaction) bool {
	//fmt.Println("Inside Local transaction verification")
	homeDir := getHomeDir()
	// get the local user balances
	userBalanceFileDir := homeDir + "/" + ACC_BALANCE_FILENAME
	localUserAccs := readLocalUserAccs(userBalanceFileDir)

	for i := 0; i < len(localUserAccs); i++ {
		trn := localUserAccs[i]

		// if the user has sufficient balance
		if transaction.BuyerId == trn.UserId && trn.FiatBalance >= transaction.BuyerPayable {
			fmt.Println("Local copy shows buyer has the balance")
			return true
		}
	}
	return false
}

//function to get a specific transaction from the central database
func getTransaction(transactionId string) Transaction {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := db.Transactions.Find(ctx, bson.M{"tId": transactionId})
	if err != nil {
		log.Fatal(err)
		fmt.Println("Failed to get the transaction data from db")
	}

	var transaction []Transaction
	if err = cursor.All(ctx, &transaction); err != nil {
		fmt.Println("Failed to write the transaction l:139 (localBlockchain.go)")
		log.Fatal(err)
	}
	return transaction[0]

}

// function to let validators check whether the latest hash block matches the one in the database or not
func checkBlock(latestCentralBlock Block) bool {

	// get the latest block in the central database
	homeDir := getHomeDir()
	// get the local blockchain
	blockchainFileDir := homeDir + "/" + LOCAL_BLOCKCHAIN
	localBlockchain := readLocalBlockchain(blockchainFileDir)

	// loop through the local blockchain
	for i := 0; i < len(localBlockchain); i++ {
		localBlock := localBlockchain[i]
		// find the db block in the local copy

		if localBlock.Hash == latestCentralBlock.Hash {
			fmt.Println("Match found\n")

			// use the nonce of the db block and hash the local block and compare this hash with the db one
			//nonce := latestCentralBlock.Nonce
			//localBlock.Nonce = nonce

			fmt.Println("Local block with db nonce\n", localBlock)
			// assign the nonce of the central blockchain to the local one and conpare the hashses
			hash := calculateHash(localBlock)
			fmt.Println("Hash of local one", hash)
			fmt.Println("Hash from db", latestCentralBlock.Hash)

			// local block's hash matches with the central one
			if hash == latestCentralBlock.Hash {
				fmt.Println("Local and db hash matched")
				return true
			}
		} else {
			fmt.Println("Match not found")
		}
	}

	fmt.Println("Local and db hash did not match")
	return false

}

//funciton to update the list of validators who checked the block with no repeat validators
func updateCheckedValidators(hash string) {
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
	blockValidators := data.Validators
	fmt.Println("The validators that checked", blockValidators)
	for i := 0; i < len(blockValidators); i++ {
		vId := blockValidators[i]
		fmt.Println("Checked validator", vId, "  Logged in user", loggedInUser)
		if loggedInUser == vId {
			found = true
		} else if i == len(blockValidators) {
			// checked all values and not found validator
			found = false
			// then we add the validator
			blockValidators = append(blockValidators, loggedInUser)
		}
	}

	// if the current validator is not part of the list, then update the list
	if !found {

		_, err := db.BlockInfo.UpdateOne(
			ctx,
			bson.M{"hash": hash},
			bson.D{
				{"$push", bson.D{{"validators", loggedInUser}}},
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

//function to get the metadata of the block given the hash of the block
func getBlockMetadata(hash string) BlockInfo {
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
	return blockData[0]
}
