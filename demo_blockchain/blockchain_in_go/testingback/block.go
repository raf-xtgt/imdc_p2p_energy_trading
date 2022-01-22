package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
	"time"
)

const difficulty = 1

type Block struct {
	Index      int
	Data       []TransactionData
	Hash       string
	PrevHash   string
	Difficulty int
	Nonce      string
}

type TransactionData struct {
	Timestamp string
	BPM       int
}

func stringifyContent(data []TransactionData) string {
	var record string
	for i := 0; i < len(data); i++ {
		record = record + data[i].Timestamp + strconv.Itoa(data[i].BPM)
	}
	return record
}
func calculateHash(block Block) string {
	transactioninString := stringifyContent(block.Data)
	record := strconv.Itoa(block.Index) + block.PrevHash + block.Nonce + transactioninString
	h := sha256.New()
	h.Write([]byte(record))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
}

func generateBlock(oldBlock Block) Block {
	var newBlock Block

	newBlock.Index = oldBlock.Index + 1
	newBlock.Data = pending
	newBlock.PrevHash = oldBlock.Hash
	newBlock.Difficulty = difficulty

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

func generateTransactionData(BPM int) TransactionData {
	var newTransactionData TransactionData
	t := time.Now()

	newTransactionData.Timestamp = t.String()
	newTransactionData.BPM = BPM

	return newTransactionData
}

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

	return true
}

func replaceChain(newBlocks []Block) {
	if len(newBlocks) > len(Blockchain) {
		Blockchain = newBlocks
	}
}

func replacePendingChain(newTransactionData []TransactionData) {
	if len(newTransactionData) > len(pending) {
		pending = newTransactionData
	}
}

/* func movePendingtoBlockChain() {
	var validatedBlock = pending[0]
	newBlockchain := append(Blockchain, validatedBlock)
	replaceChain(newBlockchain)

	pending = append(pending[:0], pending[0+1:]...)
} */

/* func deleteFirstPending() {
	pending = append(pending[:0], pending[0+1:]...)
} */

func isHashValid(hash string, difficulty int) bool {
	prefix := strings.Repeat("0", difficulty)
	return strings.HasPrefix(hash, prefix)
}

/* func validateBlock() {
	checkHash := calculateHash(pending[0])
	if checkHash == pending[0].Hash && isBlockValid(pending[0], Blockchain[len(Blockchain)-1]) {
		movePendingtoBlockChain()
	} else {
		deleteFirstPending()
	}
	writejson()
}
*/
var Blockchain []Block
var pending []TransactionData
