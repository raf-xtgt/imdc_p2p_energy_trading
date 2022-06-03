package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

/** This script is run only to initialize the first few blocks **/
const difficulty = 1
const firstBlock = "Genesis"
const normalBlock = "Normal"

// create the genesis block
func createGenesisBlock(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Methods", "GET,POST,OPTIONS,DELETE,PUT")

	// genesis bid
	var genBid Bid
	genBid.SellerId = "genesisBid"
	genBid.OptEnFromSeller = 0.00
	genBid.OptSellerReceivable = 0.00
	genBid.SellerFiatBalance = 0.00
	genBid.SellerEnergyBalance = 0.00

	// genesis transaction
	var transaction Transaction
	transaction.BuyerId = "genesis"
	transaction.BuyerPayable = 0.00
	transaction.BuyerEnReceivableFromAuction = 0.00
	transaction.BuyerEnReceivableFromTNB = 0.00

	// add the genesis bid to the auctinBids array
	var genesisBids []Bid
	genesisBidsAlt := append(genesisBids, genBid)
	transaction.AuctionBids = genesisBidsAlt
	transaction.TNBReceivable = 0.00

	var genesis Block
	genesis.Index = -1
	var genesisTransactions []Transaction
	genesisTranAlt := append(genesisTransactions, transaction)
	genesis.Data = genesisTranAlt
	genesis.PrevHash = ""
	genesis.Hash = ""

	// hash and mine the block
	generateBlock(genesis, firstBlock, genesisTransactions)
}

//get existing blocks
func updateChain(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Methods", "GET,POST,OPTIONS,DELETE,PUT")

	// get the updated local copy of blockchain and user accountsfirst
	createLocalCopies()

	// if there is a new block then verify that first
	//verifyCentralBlockchain()
	trigger := checkForNewBlocks()
	if trigger.NewBlockExists {
		verifyCentralBlockchain()
	}

	latestBlock := getLatestBlock()
	fmt.Println("The latest block is", latestBlock)

	// get a list of unverified transactions
	var newTransactionPool []Transaction = getTransactionPool()
	//fmt.Println("The three transactions", newTransactionPool)
	//validator1
	var blockTransactions []Transaction

	for j := 0; j < len(newTransactionPool); j++ {

		if verifyTransaction(newTransactionPool[j]) {
			// set transaction as verified
			//fmt.Println("Transaction Verified!!!")
			newTransaction := updateTransactionVerification(newTransactionPool[j])
			blockTransactions = append(blockTransactions, newTransaction)
			// blockTransactions = append(blockTransactions, newTransaction)
			// add the pool to a block.

		}

	}
	// if atleast one transaction is valid, then we add it to the block
	if len(blockTransactions) > 0 {
		generateBlock(latestBlock, normalBlock, blockTransactions)
	}

}

// check if any new blocks are made by other validators by checking the "trigger" collection
func checkForNewBlocks() Trigger {
	// to prevent timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	//fmt.Println("yo")
	cursor, err := db.Trigger.Find(ctx, bson.M{})
	if err != nil {
		fmt.Println("Failed to get trigger l:118 (blockchain.go)")
	}

	var newBlockExists []Trigger
	if err = cursor.All(ctx, &newBlockExists); err != nil {
		fmt.Println("Failed to write l:123 (blockchain.go)")
		log.Fatal(err)
	}
	return newBlockExists[0]

}

func getCurrentBlockchain() []Block {
	// to prevent timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// get all the blocks from the database

	//cursor, err := db.Blockchain.Find(ctx, bson.M{})
	cursor, err := db.Blockchain.Find(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
		fmt.Println("Failed to load blockchain from db")
	}

	var blockchain []Block
	if err = cursor.All(ctx, &blockchain); err != nil {
		log.Fatal(err)
		fmt.Println("Failed to load blockchain in list")
	}

	return blockchain
}

// funciton to send the current blockchain to the frontend
func sendBlockchainToFrontend(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Methods", "GET,POST,OPTIONS,DELETE,PUT")

	var response BlockchainResponse
	currentChain := getCurrentBlockchain()
	response.Blockchain = currentChain

	respondWithJSON(w, r, http.StatusCreated, response)
	return
}

type LatestIndex struct {
	Index       int "json: index"
	NewBlockNum int "json:newBlockNum"
}

// get the last mined block
func getLatestBlock() Block {
	// first get latest index

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

	// find the block with the latest index
	cursor2, err := db.Blockchain.Find(ctx, bson.M{"index": latestIndex[0].Index})
	if err != nil {
		log.Fatal(err)
	}
	var latestBlock []Block
	if err = cursor2.All(ctx, &latestBlock); err != nil {
		log.Fatal(err)
	}

	return latestBlock[0]

}

// get a pool of transactions from database
func getTransactionPool() []Transaction {
	// to prevent timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// get the transactions that have not been verified yet
	cursor, err := db.Transactions.Find(ctx, bson.M{"verified": false})
	if err != nil {
		log.Fatal(err)
	}

	var transactions []Transaction
	if err = cursor.All(ctx, &transactions); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Inside getTransactionPool", len(transactions))
	return transactions

}

// to verify whether the buyer of the transaction has sufficient balance or not
func verifyTransaction(transaction Transaction) bool {
	// to prevent timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// find the relevant data
	cursor, err := db.UserAccBalance.Find(ctx, bson.M{"userid": transaction.BuyerId})
	if err != nil {
		log.Fatal(err)
	}
	// fmt.Println("Inside verification function")
	// fmt.Println("The transaction is", transaction)

	var userAccount []AccountBalance
	if err = cursor.All(ctx, &userAccount); err != nil {
		log.Fatal(err)
		fmt.Println("User Account not found")
		return false
	} else {
		fmt.Println(userAccount)
		if userAccount[0].FiatBalance >= transaction.BuyerPayable {
			fmt.Println("User: ", transaction.BuyerId, "has sufficient account balance")
			return true
		} else {
			fmt.Println("User: ", transaction.BuyerId, "does not have sufficient account balance")
			return false
		}
	}
}

// set isverified to true in db
func updateTransactionVerification(transaction Transaction) Transaction {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// set the isVerified to true
	_, err := db.Transactions.UpdateOne(
		ctx,
		bson.M{"tId": transaction.TId},
		bson.D{
			{"$set", bson.D{{"verified", true}}},
		},
	)
	incrementChecks(transaction.TId)

	// if the update fails
	if err != nil {
		log.Fatal(err)
		fmt.Println("Failed to update the transaction")
	}
	transaction.Verified = true

	return transaction

}

// function to increment the number of validator checks in a transaction in the central db
func incrementChecks(transactionId string) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// increment checks
	_, err := db.Transactions.UpdateOne(
		ctx,
		bson.M{"tId": transactionId},
		bson.D{
			{"$inc", bson.D{{"checks", 1}}},
		},
	)

	// if the update fails
	if err != nil {
		log.Fatal(err)
		fmt.Println("Failed to update the transaction checks")
	}

	fmt.Println("Validator successfully incremented checks")

}

// calculate nonce, and the hash of the block
func generateBlock(oldBlock Block, blockType string, transactions []Transaction) {
	var newBlock Block

	newBlock.Index = oldBlock.Index + 1
	newBlock.PrevHash = oldBlock.Hash

	// for genesis block only
	if blockType == firstBlock {
		newBlock.Data = oldBlock.Data
	} else {
		// for all other blocks
		newBlock.Data = transactions
	}

	fmt.Println("New block", newBlock)

	// time how long it takes to mine
	loc, _ := time.LoadLocation("UTC")
	startTime := time.Now().In(loc)
	timeStr := startTime.String()
	fmt.Println("Start time before mining starts", timeStr)

	// This one is the mining algorithm to find the nonce that suits the hash requirement (how many 0)
	// Given in mycoralhealth website
	for i := 0; ; i++ {
		hex := fmt.Sprintf("%x", i)
		newBlock.Nonce = hex
		if !isHashValid(calculateHash(newBlock), difficulty) {
			fmt.Println(calculateHash(newBlock), " incorrect hash!")
			time.Sleep(time.Second)
			continue
		} else {
			fmt.Println(calculateHash(newBlock), " correct hash found!")
			newBlock.Hash = calculateHash(newBlock)
			// add the block to database
			addBlock(newBlock)
			break
		}

	}

	// end time

	endTime := time.Now().In(loc)
	endTimeStr := endTime.String()
	fmt.Println("End time when mining stops", endTimeStr)
}

// To calculate the hash of a block
// Given in the mycoralhealth website and from code given by Teoh
func calculateHash(block Block) string {
	transactioninString := stringifyTransaction(block.Data)
	//fmt.Println("Inside calculate hash")
	record := strconv.Itoa(block.Index) +
		block.PrevHash + block.Nonce + transactioninString

	h := sha256.New()
	h.Write([]byte(record))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
}

// Check whether hash has enough 0 in front
// Code given by Teoh
func isHashValid(hash string, difficulty int) bool {
	prefix := strings.Repeat("0", difficulty)
	return strings.HasPrefix(hash, prefix)
}

// Turn everything in Transaction Data into string to calculate the hash afterwards
func stringifyTransaction(data []Transaction) string {
	//fmt.Println("Inside stringify function")

	var record string
	for i := 0; i < len(data); i++ {

		bids := data[i].AuctionBids
		var bidsStr string = ""
		for j := 0; j < len(bids); j++ {

			bidsStr = bidsStr + bids[j].SellerId + fmt.Sprintf("%v", bids[j].OptEnFromSeller) +
				fmt.Sprintf("%v", bids[j].OptSellerReceivable) +
				fmt.Sprintf("%v", bids[j].SellerFiatBalance) +
				fmt.Sprintf("%v", bids[j].SellerEnergyBalance)

		}
		record = record + data[i].BuyerId + fmt.Sprintf("%v", data[i].BuyerPayable) +
			fmt.Sprintf("%v", data[i].BuyerEnReceivableFromAuction) +
			fmt.Sprintf("%v", data[i].BuyerEnReceivableFromTNB) +
			fmt.Sprintf("%v", data[i].TNBReceivable) + bidsStr

	}
	return record
}

// add the block to the database
func addBlock(newBlock Block) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	writeBlock, err := db.Blockchain.InsertOne(ctx, newBlock)
	if err != nil {
		log.Fatal(err)
		fmt.Println("failed to write block")
		return
	}
	fmt.Println("Successfully added block", writeBlock.InsertedID)

	// add the validator who made the block as one of the "validators" who checked the block
	var blockMetadata BlockInfo
	var blockIdStr string = fmt.Sprintf("%v", writeBlock.InsertedID)
	blockMetadata.BlockId = blockIdStr[10 : len(blockIdStr)-2]
	var validators []string
	validators = append(validators, loggedInUser)
	blockMetadata.Validators = validators
	blockMetadata.Hash = newBlock.Hash
	var clerks []string
	clerks = append(clerks, "")
	blockMetadata.Clerks = clerks

	//add this block metadata in the db
	writeBlockInfo, err := db.BlockInfo.InsertOne(ctx, blockMetadata)
	if err != nil {
		log.Fatal(err)
		fmt.Println("failed to write block")
		return
	}
	fmt.Println("Successfully added block metadata", writeBlockInfo.InsertedID)

	// add the blockInfo(metadata) in the db
	updateLatestBlockIndex(newBlock.Index)

}

//function to update the blockinfo. Here we are storing the number of validators who have checked the block)

// increment the latest index to the value of the new blocks
func updateLatestBlockIndex(index int) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// set the new index
	_, err := db.LatestIndex.UpdateOne(
		ctx,
		bson.M{},
		bson.D{
			{"$inc", bson.D{{"index", 1}}},
		},
	)

	// if the update fails
	if err != nil {
		log.Fatal(err)
		fmt.Println("Failed to update the latest index")
	} else {
		fmt.Println("Successfully updated the latest index")
	}

	// increment the number of new blocks
	_, err2 := db.LatestIndex.UpdateOne(
		ctx,
		bson.M{},
		bson.D{
			{"$inc", bson.D{{"newBlockNum", 1}}},
		},
	)

	// if the update fails
	if err2 != nil {
		log.Fatal(err2)
		fmt.Println("Failed to increment number of new blocks")
	} else {
		fmt.Println("Successfully updated the number of new blocks")
	}

	// make a trigger that a new block has been made
	setTrigger(true)

	// create a local copy to include the new block
	createLocalCopies()
}

// function to update the trigger collection
func setTrigger(value bool) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// set the isVerified to true
	_, err := db.Trigger.UpdateOne(
		ctx,
		bson.M{},
		bson.D{
			{"$set", bson.D{{"newBlockExists", value}}},
		},
	)

	// if the update fails
	if err != nil {
		log.Fatal(err)
		fmt.Println("Failed to update the trigger")
	} else {
		fmt.Println("Successfully updated the trigger")
	}
}

// function to get all the userAccounts
func getAllAccs() []AccountBalance {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := db.UserAccBalance.Find(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}

	var allUserAccs []AccountBalance
	if err = cursor.All(ctx, &allUserAccs); err != nil {
		log.Fatal(err)
	}
	return allUserAccs
}

// function to find the local directory to store the local blockchain and user account data
func getHomeDir() string {
	// the default directory to store the blockchain
	dirname, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("dirname", dirname)
	return dirname
}

// create a local copy of the blockchain and the user accounts as it exists in the central db as of now
func createLocalCopies() {
	// get the path to store the local copies
	dirname := getHomeDir()

	// get the current blockchain from the central db
	currentBlockchain := getCurrentBlockchain()
	// create the local file and get the full file path
	fileDir := createLocalBlockchainFile(dirname)
	// write the blockchain to the local file
	writeLocalBlockchain(currentBlockchain, fileDir)
	//readLocalBlockchain(fileDir)

	// get the user accounts from the databse
	allAccs := getAllAccs()
	//create a local account file
	accPath := createUserAccountCopy(dirname)
	// write the local account copies
	writeToUserLocalAcc(allAccs, accPath)

}

// function to discard a block
func discardBlock(latestCentralBlock Block) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	//  delete the block
	result, err := db.Blockchain.DeleteOne(ctx, bson.M{"hash": latestCentralBlock.Hash})
	if err != nil {
		fmt.Println("Failed to discard the block blockchain.go ln:528")
		log.Fatal(err)
	}
	fmt.Printf("Discarded block %v document(s)\n", result.DeletedCount)

	// decrement the index to point to the block before the discarded one
	_, err2 := db.LatestIndex.UpdateOne(
		ctx,
		bson.M{},
		bson.D{
			{"$inc", bson.D{{"index", -1}}},
		},
	)

	// if the update fails
	if err2 != nil {
		log.Fatal(err2)
		fmt.Println("Failed to update the transaction checks")
	}

	fmt.Println("Validator successfully incremented checks")

}
