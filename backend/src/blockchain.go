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

/** This script is run only to initialize the first few blocks **/
const difficulty = 2
const firstBlock = "Genesis"
const buyer = "buyer"
const seller = "seller"

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
	generateBlock(genesis, firstBlock)
	return
}

//get existing blocks
func updateChain(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Methods", "GET,POST,OPTIONS,DELETE,PUT")
	currentBlockchain := getCurrentBlockchain()
	fmt.Println("Current blockchain")
	fmt.Println("Current blockchain", currentBlockchain)
}

func getCurrentBlockchain() []Block {
	// to prevent timeout
	mongoparams.ctx, mongoparams.cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer mongoparams.cancel()

	// get all the blocks from the database
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

// calculate nonce, and the hash of the block
func generateBlock(oldBlock Block, blockType string) {
	var newBlock Block

	newBlock.Index = oldBlock.Index + 1
	newBlock.PrevHash = oldBlock.Hash

	// for genesis block only
	if blockType == firstBlock {
		newBlock.Data = oldBlock.Data
	} else {
		// for all other blocks
		var blockTransactions [3]Transaction
		blockTransactions = getTransactions()
		fmt.Println("The pool of three transactions")
		for i := 0; i < len(blockTransactions); i++ {
			fmt.Println("BuyerId", blockTransactions[i].BuyerId)
			fmt.Println("BuyerPayable", blockTransactions[i].BuyerPayable)
		}
	}

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
			break
		}

	}

	// add the block to database
	addBlock(newBlock)
}

// get all the transactions from the database
func getTransactions() [3]Transaction {
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

	var blockTrans [3]Transaction
	blockTrans[0] = transactions[0]
	blockTrans[1] = transactions[1]
	blockTrans[2] = transactions[2]

	return blockTrans
}

// To calculate the hash of a block
// Given in the mycoralhealth website and from code given by Teoh
func calculateHash(block Block) string {
	transactioninString := stringifyTransaction(block.Data)

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
	var record string
	for i := 0; i < len(data); i++ {

		bids := data[i].AuctionBids
		var bidsStr string = ""
		for j := 0; j < len(bids); j++ {

			bidsStr = bidsStr + bids[i].SellerId + fmt.Sprintf("%v", bids[i].OptEnFromSeller) +
				fmt.Sprintf("%v", bids[i].OptSellerReceivable) +
				fmt.Sprintf("%v", bids[i].SellerFiatBalance) +
				fmt.Sprintf("%v", bids[i].SellerEnergyBalance)

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
	}
	fmt.Println("Successfully added block", writeBlock.InsertedID)

}

//verify transaction by checking whether buyer has the required amount of fiat money
/**
func verifyTransaction(transaction Transaction) bool {
	// to prevent timeout
	mongoparams.ctx, mongoparams.cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer mongoparams.cancel()

	// find the relevant data
	cursor, err := db.UserAccBalance.Find(ctx, bson.M{"buyerId": transaction.BuyerId})
	if err != nil {
		log.Fatal(err)
	}

	var userAccount AccountBalance
	if err = cursor.All(ctx, &userAccount); err != nil {
		log.Fatal(err)
		fmt.Println("User Account not found")
		return false
	} else {
		if userAccount.FiatBalance >= transaction.BuyerPayable {
			fmt.Println("User: ", transaction.BuyerId, "has sufficient account balance")
			return true
		} else {
			fmt.Println("User: ", transaction.BuyerId, "does not have sufficient account balance")
			return false
		}
	}
}

// update the fiat balance for verified transaction
func updateUserAccBalances(userAccount AccountBalance, userType string, payable float64, userId string) {
	var newBalance float64
	if userType == buyer {
		newBalance = userAccount.FiatBalance - payable
		userAccount.FiatBalance = newBalance

		// update the database
		_, err := db.UserAccBalance.UpdateOne(
			mongoparams.ctx,
			bson.M{"userid": userId},
			bson.D{
				{"$set", bson.D{{"fiatbalance", payable}}},
			},
		)

		// if the update fails
		if err != nil {
			log.Fatal(err)
			return false
		}

	}
}

**/
