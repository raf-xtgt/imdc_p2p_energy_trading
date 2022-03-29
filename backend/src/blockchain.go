package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

type OrderedBlocks []Block

func (b OrderedBlocks) Len() int {
	return len(b)
}

func (b OrderedBlocks) Less(i, j int) bool {
	return b[i].Index < b[j].Index
}

func (b OrderedBlocks) Swap(i, j int) {
	b[i], b[j] = b[j], b[i]
}

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
	currentBlockchain := getCurrentBlockchain()
	fmt.Println("Current blockchain", currentBlockchain)

	latestBlock := getLatestBlock()
	fmt.Println("The latest block is", latestBlock)

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

func getCurrentBlockchain() []Block {
	// to prevent timeout
	mongoparams.ctx, mongoparams.cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer mongoparams.cancel()

	// get all the blocks from the database

	//cursor, err := db.Blockchain.Find(mongoparams.ctx, bson.M{})
	cursor, err := db.Blockchain.Find(mongoparams.ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
		fmt.Println("Failed to load blockchain from db")
	}

	var blockchain []Block
	if err = cursor.All(mongoparams.ctx, &blockchain); err != nil {
		log.Fatal(err)
		fmt.Println("Failed to load blockchain in list")
	}

	return blockchain
}

type LatestIndex struct {
	Index int "json: index"
}

// get the last mined block
func getLatestBlock() Block {
	// first get latest index

	mongoparams.ctx, mongoparams.cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer mongoparams.cancel()

	// get the latest index. The collection always has only one document
	cursor, err := db.LatestIndex.Find(mongoparams.ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	var latestIndex []LatestIndex
	if err = cursor.All(mongoparams.ctx, &latestIndex); err != nil {
		log.Fatal(err)
	}

	// find the block with the latest index
	cursor2, err := db.Blockchain.Find(mongoparams.ctx, bson.M{"index": latestIndex[0].Index})
	if err != nil {
		log.Fatal(err)
	}
	var latestBlock []Block
	if err = cursor2.All(mongoparams.ctx, &latestBlock); err != nil {
		log.Fatal(err)
	}

	return latestBlock[0]

}

// get a pool of transactions from database
func getTransactionPool() []Transaction {
	// to prevent timeout
	mongoparams.ctx, mongoparams.cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer mongoparams.cancel()

	// get the transactions that have not been verified yet
	cursor, err := db.Transactions.Find(mongoparams.ctx, bson.M{"verified": false})
	if err != nil {
		log.Fatal(err)
	}

	var transactions []Transaction
	if err = cursor.All(mongoparams.ctx, &transactions); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Inside getTransactionPool", len(transactions))
	return transactions

}

// to verify whether the buyer of the transaction has sufficient balance or not
func verifyTransaction(transaction Transaction) bool {
	// to prevent timeout
	mongoparams.ctx, mongoparams.cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer mongoparams.cancel()

	// find the relevant data
	cursor, err := db.UserAccBalance.Find(mongoparams.ctx, bson.M{"userid": transaction.BuyerId})
	if err != nil {
		log.Fatal(err)
	}
	// fmt.Println("Inside verification function")
	// fmt.Println("The transaction is", transaction)

	var userAccount []AccountBalance
	if err = cursor.All(mongoparams.ctx, &userAccount); err != nil {
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
	mongoparams.ctx, mongoparams.cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer mongoparams.cancel()

	// set the isVerified to true
	_, err := db.Transactions.UpdateOne(
		mongoparams.ctx,
		bson.M{"tId": transaction.TId},
		bson.D{
			{"$set", bson.D{{"verified", true}}}, //FOR NOW
		},
	)

	// if the update fails
	if err != nil {
		log.Fatal(err)
		fmt.Println("Failed to update the transaction")
	}
	transaction.Verified = true
	return transaction

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

}

// To calculate the hash of a block
// Given in the mycoralhealth website and from code given by Teoh
func calculateHash(block Block) string {
	transactioninString := stringifyTransaction(block.Data)
	fmt.Println("Inside calculate hash")
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
	fmt.Println("Inside stringify function")

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
	mongoparams.ctx, mongoparams.cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer mongoparams.cancel()
	writeBlock, err := db.Blockchain.InsertOne(mongoparams.ctx, newBlock)
	if err != nil {
		log.Fatal(err)
		fmt.Println("failed to write block")
		return
	}
	fmt.Println("Successfully added block", writeBlock.InsertedID)

	updateLatestBlockIndex(newBlock.Index)

}

// increment the latest index to the value of the new blocks
func updateLatestBlockIndex(index int) {
	mongoparams.ctx, mongoparams.cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer mongoparams.cancel()

	// set the isVerified to true
	_, err := db.LatestIndex.UpdateOne(
		mongoparams.ctx,
		bson.M{},
		bson.D{
			{"$set", bson.D{{"index", index}}},
		},
	)

	// if the update fails
	if err != nil {
		log.Fatal(err)
		fmt.Println("Failed to update the latest index")
	} else {
		fmt.Println("Successfully updated the latest index")
	}

}
