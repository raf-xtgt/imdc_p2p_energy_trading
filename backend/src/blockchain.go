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
)

/** This script is run only to initialize the first few blocks **/
const difficulty = 2
const firstBlock = "Genesis"

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

// calculate nonce, and the hash of the block
func generateBlock(oldBlock Block, blockType string) {
	var newBlock Block

	newBlock.Index = oldBlock.Index + 1
	newBlock.PrevHash = oldBlock.Hash

	// for genesis block only
	if blockType == firstBlock {
		newBlock.Data = oldBlock.Data
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
