package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"
)

// Role:
// 1 -> Admin (can transfer money without limit + Have Private/Public key)
// 2 -> Leader (Have Private/Public key to encrypt block)
// 3 -> Validator (Can store blockchain)
// 4 -> Normal (Nothing)

//Missing: Password and other user info
// Userid data type can be changed
type User struct {
	Userid         int
	Username       string
	Pubkey         []byte
	Privkey        []byte
	Role           int
	Blockchain     []Block
	BlockchainHash string
}

// Create a new admin
// The id, username for the admin can be changed
func CreateAdmin() {
	newid := 99
	var newUser User
	newUser.Userid = newid
	newUser.Username = "Admin"
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic(err)
	}
	publicKey := privateKey.PublicKey

	newUser.Privkey = x509.MarshalPKCS1PrivateKey(privateKey)

	newUser.Pubkey = x509.MarshalPKCS1PublicKey(&publicKey)

	newUser.Role = 1

	// Add the admin to the user list (profiles)
	newProfiles := append(profiles, newUser)
	profiles = newProfiles

	writeUserjson()

}

// Create a new leader
func CreateLeader() {
	newid := len(profiles)
	var newUser User
	newUser.Userid = newid
	newUser.Username = "Leader"
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic(err)
	}
	publicKey := privateKey.PublicKey

	newUser.Privkey = x509.MarshalPKCS1PrivateKey(privateKey)

	newUser.Pubkey = x509.MarshalPKCS1PublicKey(&publicKey)

	newUser.Role = 2

	newProfiles := append(profiles, newUser)
	profiles = newProfiles

	writeUserjson()
}

// Create a validator
func CreateValidator() {
	newid := len(profiles)
	var newUser User
	newUser.Userid = newid
	newUser.Username = "Validator"

	newUser.Role = 3

	if Blockchain == nil {
		genesisBlock := createGenesisBlock()
		newUser.Blockchain = append(newUser.Blockchain, genesisBlock)
	} else {
		newUser.Blockchain = Blockchain
	}
	newProfiles := append(profiles, newUser)
	profiles = newProfiles

	createTopUp(newUser.Userid)

	writeUserjson()
}

// Create a new user
func CreateNewUser() {
	newid := len(profiles)
	var newUser User
	newUser.Userid = newid
	newUser.Username = ""

	newUser.Role = 4

	newProfiles := append(profiles, newUser)
	profiles = newProfiles

	createTopUp(newUser.Userid)

	writeUserjson()
}

// create a transaction data block that represents the user receive some money
func createTopUp(userid int) {
	var topup TransactionData
	topup.Buyer = 99 //admin, change this dynamically if the admin user id is going to change
	topup.Seller = userid
	topup.Money = 100
	topup.Energy = 0
	t := time.Now()
	topup.Timestamp = t.String()

	newPending := append(pending, topup)
	pending = newPending

}

// Just write user list (profiles) into json file same as writejson and readjson
func readUserjson() {
	jsonFile, err := os.Open("users.json")
	if err != nil {
		log.Fatal(err)
	}
	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &profiles)
}

func writeUserjson() {
	jsonFile, _ := json.MarshalIndent(profiles, "", " ")

	err := ioutil.WriteFile("users.json", jsonFile, 0644)
	if err != nil {
		panic(err)
	}
}

// Obtain the user balance
func GetUserBalance(userid int) int {
	var trackfrom int
	// Exception for the first time, because we dont have balance block yet
	if len(Blockchain) < 5 {
		trackfrom = 1
	} else {
		// The balance block is appended every 5th block, "/" is quotient without returning remainder
		trackfrom = (len(Blockchain) / 5) * 5
	}

	var balance int
	// Start from the nearest balance block
	for i := (trackfrom - 1); i < len(Blockchain); i++ {
		for j := 0; j < len(Blockchain[i].Data); j++ {
			// If user sells energy, they received money, hence balance += money
			if Blockchain[i].Data[j].Seller == userid {
				balance = balance + Blockchain[i].Data[j].Money
			}
			// If user buys energy, they spends money, hence balance -= money
			if Blockchain[i].Data[j].Buyer == userid {
				balance = balance - Blockchain[i].Data[j].Money
			}
		}
	}
	// Calculate the balance based on the transaction data in pending as well
	for i := 0; i < len(pending); i++ {
		if pending[i].Seller == userid {
			balance = balance + pending[i].Money
		}
		if pending[i].Buyer == userid {
			balance = balance - pending[i].Money
		}
	}
	return balance
}

//create balance block
func GenerateBalanceBlock() {
	for i := 0; i < len(profiles); i++ {
		// Currently, only validators and normal users can make a trade
		if profiles[i].Role == 4 || profiles[i].Role == 3 {
			uid := profiles[i].Userid
			balance := GetUserBalance(uid)
			var topup TransactionData
			topup.Buyer = 99 //admin, change this if the admin id is changed
			topup.Seller = uid
			topup.Money = balance
			topup.Energy = 0
			t := time.Now()
			topup.Timestamp = t.String()

			newPending := append(pending, topup)
			pending = newPending
		}
	}
	newBalanceBlock := generateBlock(Blockchain[len(Blockchain)-1])

	SendBlock(newBalanceBlock)
	writejson()

	// Clear and reset the pending array
	clearPending()
}

// Check which blockchain hash is the majority
func CheckOccurenceBlockchain() string {
	var BlockchainHashSet []string
	var count []int
	for i := 0; i < len(profiles); i++ {
		if profiles[i].Role == 3 {
			// Check whether the blockchain hash exists in the array
			Index, Exist := checkHashExist(profiles[i].BlockchainHash, BlockchainHashSet)
			// If exist, its correspond count + 1
			if Exist {
				count[Index] = count[Index] + 1
			} else {
				// New Hash, append to the current array
				BlockchainHashSet = append(BlockchainHashSet, profiles[i].BlockchainHash)
				count = append(count, 1)
			}

		}
	}
	fmt.Println(BlockchainHashSet)
	fmt.Println(count)
	//Find max/majority
	max := count[0]
	var MaxOccurenceHash string
	for i := 0; i < len(count); i++ {
		if count[i] > max {
			max = count[i]
		}
		MaxOccurenceHash = BlockchainHashSet[i]
	}
	return MaxOccurenceHash

}

// Simply check whether a Hash exist in the BlockchainHashSet array
func checkHashExist(Hash string, Array []string) (int, bool) {
	for i, j := range Array {
		if j == Hash {
			return i, true
		}
	}
	return -1, false
}

var profiles []User // userlist
