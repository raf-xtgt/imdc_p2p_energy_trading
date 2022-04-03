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
		fmt.Println("File created successfully", filename)
		defer file.Close()
	} else {
		fmt.Println("File already exists!", filename)

	}

	return filename

}

// write to the blockchain in the file
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

	latestCentralBlock := getLatestBlock()
	blockTransactions := latestCentralBlock.Data

	for j := 0; j < len(blockTransactions); j++ {
		transaction := blockTransactions[j]
		if localTrnVerification(transaction) {
			// increment validator check in the central database
			incrementChecks(transaction.TId)

			// get the recently incremented transaction from the database
			updatedTrn := getTransaction(transaction.TId)

			// if all validators have checked that the user has sufficient balance
			if updatedTrn.Checks >= TOTAL_VALIDATORS {
				// use the nonce of the latest block and check whether its hash matches or not
				return checkBlock(latestCentralBlock)
			} else {
				fmt.Println("Both validators have not verified the block")
			}

		} else {
			fmt.Println("Local transaction verification failed")
		}

	}
	return false

}

// other not currently logged in validators check whether user has sufficient balance or not
func localTrnVerification(transaction Transaction) bool {
	fmt.Println("Inside Local transaction verification")
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
	mongoparams.ctx, mongoparams.cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer mongoparams.cancel()

	cursor, err := db.Transactions.Find(mongoparams.ctx, bson.M{"tId": transactionId})
	if err != nil {
		log.Fatal(err)
		fmt.Println("Failed to get the transaction data from db")
	}

	var transaction []Transaction
	if err = cursor.All(mongoparams.ctx, &transaction); err != nil {
		fmt.Println("Failed to write the transaction l:139 (localBlockchain.go)")
		log.Fatal(err)
	}
	return transaction[0]

}

// function to let not currently logged in validators check whether the latest hash block matches to the one in the database or not
func checkBlock(latestCentralBlock Block) bool {

	// get the latest block in the central database
	homeDir := getHomeDir()
	// get the local blockchain
	blockchainFileDir := homeDir + "/" + LOCAL_BLOCKCHAIN
	localBlockchain := readLocalBlockchain(blockchainFileDir)
	latestLocalBlock := localBlockchain[len(localBlockchain)-1] // the last block is the latest and new block

	nonce := latestCentralBlock.Nonce
	latestLocalBlock.Nonce = nonce

	// assign the nonce of the central blockchain to the local one and conpare the hashses
	hash := calculateHash(latestLocalBlock)

	// local block's hash matches with the central one
	if hash == latestCentralBlock.Hash {
		fmt.Println("Local and db hash matched")
		return true
	} else {
		fmt.Println("Local and db hash did not match")
		return false
	}

}
