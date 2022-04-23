package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
	"time"
)

const difficulty = 0 // the number of 0's in the Hash

// Maybe can remove Difficulty
type Block struct {
	Index      int
	Data       []TransactionData
	Hash       string
	PrevHash   string
	Difficulty int
	Nonce      string
	LeaderID   int
}

// The way to read it: (Buyer) spends (Money) to buy (Energy) from (Seller)
// Buyer (- Money); Seller (+ Money)
type TransactionData struct {
	Timestamp string
	Buyer     int
	Seller    int
	Money     int
	Energy    int
}

// This AltTransactionData is created because the HTTP request I received is all string,
// You can change the way to handle the HTTP request as u want
// And also I didn't include the time in my front end as well
type AltTransactionData struct {
	Buyer  string
	Seller string
	Money  string
	Energy string
}

// To obtain the hash of the blockchain
func calculateBlockchainHash(chain []Block) string {

	// Turn everything inside the blockchain to string
	var record string
	for i := 0; i < len(chain); i++ {
		transactionString := stringifyTransaction(chain[i].Data)
		record = record + strconv.Itoa(chain[i].Index) + strconv.Itoa(chain[i].LeaderID) + chain[i].PrevHash + chain[i].Nonce + transactionString
	}

	h := sha256.New()
	h.Write([]byte(record))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
}

// Turn everything in Transaction Data into string to calculate the hash afterwards
func stringifyTransaction(data []TransactionData) string {
	var record string
	for i := 0; i < len(data); i++ {
		record = record + data[i].Timestamp + strconv.Itoa(data[i].Buyer) +
			strconv.Itoa(data[i].Seller) + strconv.Itoa(data[i].Money) + strconv.Itoa(data[i].Energy)
	}
	return record
}

// To calculate the hash of a block
// Given in the mycoralhealth website
func calculateHash(block Block) string {
	transactioninString := stringifyTransaction(block.Data)

	record := strconv.Itoa(block.Index) + strconv.Itoa(block.LeaderID) +
		block.PrevHash + block.Nonce + transactioninString

	h := sha256.New()
	h.Write([]byte(record))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
}

// To generate Block from Transaction data in pending
func generateBlock(oldBlock Block) Block {
	var newBlock Block

	newBlock.Index = oldBlock.Index + 1
	newBlock.Data = pending
	newBlock.PrevHash = oldBlock.Hash
	newBlock.Difficulty = difficulty
	newBlock.LeaderID = NextLeader() // Always the next leader from the previous one

	// This one is the mining algorithm to find the nonce that suits the hash requirement (how many 0)
	// Given in mycoralhealth website
	for i := 0; ; i++ {
		hex := fmt.Sprintf("%x", i)
		newBlock.Nonce = hex
		if !isHashValid(calculateHash(newBlock), newBlock.Difficulty) {
			fmt.Println(calculateHash(newBlock), " do more work!")
			time.Sleep(time.Second)
			continue
		} else {
			fmt.Println(calculateHash(newBlock), " work done!")
			newBlock.Hash = calculateHash(newBlock)
			break
		}

	}
	return newBlock
}

// Check whether the block is valid
// The new block index must be old block index + 1
// The new block prevhash must be == old block hash
// The hash calculated must be same
// The new block must have some transaction data inside it
func isBlockValid(newBlock, oldBlock Block) bool {
	if oldBlock.Index+1 != newBlock.Index {
		return false
	}

	if oldBlock.Hash != newBlock.PrevHash {
		return false
	}

	if calculateHash(newBlock) != newBlock.Hash {
		return false
	}

	if len(newBlock.Data) == 0 {
		return false
	}

	return true
}

// Just replace the old blockchain with a new one
func replaceChain(newBlocks []Block) {
	if len(newBlocks) > len(Blockchain) {
		Blockchain = newBlocks
	}
}

// Validate the transaction
// Make sure the Buyer has enough balance to make the purchase
// Can add: Check Seller has enough energy to sell
// Can add: Check user id exists
// Can add: Other validating techniques
func validateTransaction(TxData TransactionData) bool {
	//check enough balance
	var validity bool

	balance := GetUserBalance(TxData.Buyer)
	if balance >= TxData.Money {
		validity = true
	} else {
		validity = false
	}

	if TxData.Buyer == 99 { // Should be admin ID
		validity = true
	}

	return validity
}

// Check whether hash has enough 0 in front
func isHashValid(hash string, difficulty int) bool {
	prefix := strings.Repeat("0", difficulty)
	return strings.HasPrefix(hash, prefix)
}

// Create a genesis block
/**
type Block struct {
	Index      int
	Data       []TransactionData
	Hash       string
	PrevHash   string
	Difficulty int
	Nonce      string
	LeaderID   int
}
**/
func createGenesisBlock() Block {
	genesisTransactionData := TransactionData{}
	genesisTransactionData = TransactionData{"0", 0, 0, 0, 0}
	var genesisTransactionDataChain []TransactionData
	genesisTransactionDataChainAlt := append(genesisTransactionDataChain, genesisTransactionData)
	genesisTransactionDataChain = genesisTransactionDataChainAlt
	genesisBlock := Block{}
	genesisBlock = Block{0, genesisTransactionDataChain, calculateHash(genesisBlock), "", difficulty, "", 0}
	return genesisBlock
}

// Check whether the main blockchain == the majority blockchain stored by the validators
func checkBlockchainValid() {

	fmt.Println(calculateBlockchainHash(Blockchain))

	//Acquire the most frequent appear blockchain hash
	MostOccurenceHash := CheckOccurenceBlockchain()

	if MostOccurenceHash == calculateBlockchainHash(Blockchain) {
		fmt.Println("Current Blockchain is Correct")
	} else {
		fmt.Println("Current Blockchain is Updated")
		// Obtain the blockchain corresponding to the blockchain hash
		for i := 0; i < len(profiles); i++ {
			if MostOccurenceHash == profiles[i].BlockchainHash {
				Blockchain = profiles[i].Blockchain
				writejson() // store the blockchain to the json file
				return
			}
		}
	}
}

// A function to find the next leader,
// if there is no leader behind the current one, will recycle to the first one
func NextLeader() int {
	var LeaderIDArray []int
	var NextLeaderID int

	// Obtain a list contain all the leaders
	for i := 0; i < len(profiles); i++ {
		if profiles[i].Role == 2 {
			LeaderIDArray = append(LeaderIDArray, profiles[i].Userid)
		}
	}

	// An exception for the first block after the genesis block
	if len(Blockchain) == 1 {
		NextLeaderID = LeaderIDArray[0]
	} else {
		// Check if the leader in the last blockchain is the last leader in the leader array
		if Blockchain[len(Blockchain)-1].LeaderID == LeaderIDArray[len(LeaderIDArray)-1] {
			// if yes, the leader for the next block will be the first leader
			NextLeaderID = LeaderIDArray[0]
		} else {
			// Just find the next leader who is responsible
			for j := 0; j < len(LeaderIDArray); j++ {
				if Blockchain[len(Blockchain)-1].LeaderID == LeaderIDArray[j] {
					NextLeaderID = LeaderIDArray[j+1]
				}
			}
		}
	}

	return NextLeaderID
}

var Blockchain []Block
var pending []TransactionData
