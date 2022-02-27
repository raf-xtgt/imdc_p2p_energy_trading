package main

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"

	"github.com/davecgh/go-spew/spew"
)

const protocol = "tcp4"

// The format to send the signed block through the TCP port
type BlockMsgBlock struct {
	BlockToSend Block
	MsgHashSum  []byte
	Signature   []byte
}

// Whole concept of this p2p module:
// A tcp port 3000 is initiated
// Block created is signed with leader's private key
// Signed block is send to the tcp 3000
// A read tcp port function is run forever to read any message in tcp 3000
// It collects the msg block and verifies the block inside using leader's public key
// The verified block is then append to the main blockchain and other validator's blockchain

// Actually the tcp port is not really needed and can simply replaced by calling a function
// Because everything is run on the server, hence very little human action involved
// But in order to implement the digital signing algorithm to show some proof of authority
// Hence a tcp port is used to introduced a secure passage using the digital signing algorithm

// Send the block to the tcp port
func SendBlock(BlockToSend Block) {
	nodeAddr := ":3000"
	conn, err := net.Dial(protocol, nodeAddr)

	if err != nil {
		log.Panic(err)
	}

	defer conn.Close()

	// Obtain leader's private key
	var LeaderPrivateKey *rsa.PrivateKey
	for i := 0; i < len(profiles); i++ {
		if profiles[i].Userid == BlockToSend.LeaderID {
			// convert the private key store in the user list (profiles) to a KEY type
			LeaderPrivateKey, err = x509.ParsePKCS1PrivateKey(profiles[i].Privkey)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
	// The digital signing algorithm is based on website:
	// https://www.sohamkamani.com/golang/rsa-encryption/
	TxBlock, _ := json.MarshalIndent(BlockToSend, "", " ")
	msgHash := sha256.New()
	_, err = msgHash.Write(TxBlock)
	if err != nil {
		panic(err)
	}
	msgHashSum := msgHash.Sum(nil)

	// Sign the transaction data
	signature, err := rsa.SignPSS(rand.Reader, LeaderPrivateKey, crypto.SHA256, msgHashSum, nil)
	if err != nil {
		panic(err)
	}

	// Create a message block to send all the information needed
	newMessageBlock := BlockMsgBlock{}
	newMessageBlock = BlockMsgBlock{BlockToSend, msgHashSum, signature}
	newMessageBlockToSend, _ := json.MarshalIndent(newMessageBlock, "", " ")
	data := []byte(newMessageBlockToSend)

	// write the message block to the tcp port
	_, err = io.Copy(conn, bytes.NewReader(data))
	if err != nil {
		log.Panic(err)
	}

}

// Read the data from tcp port
func ReadBlock(conn net.Conn) {
	var MsgReceived BlockMsgBlock
	var BlockReceived Block
	req, err := ioutil.ReadAll(conn)
	defer conn.Close()
	if err != nil {
		log.Panic(err)
	}

	// Convert the message from string to struct type
	json.Unmarshal(req, &MsgReceived)

	// Obtain Leader's public key
	var LeaderPublicKey *rsa.PublicKey
	for i := 0; i < len(profiles); i++ {
		if profiles[i].Userid == MsgReceived.BlockToSend.LeaderID {
			LeaderPublicKey, err = x509.ParsePKCS1PublicKey(profiles[i].Pubkey)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
	err = rsa.VerifyPSS(LeaderPublicKey, crypto.SHA256, MsgReceived.MsgHashSum, MsgReceived.Signature, nil)
	if err != nil {
		fmt.Println("could not verify signature: ", err)
		return
	} else if err == nil {
		fmt.Println("Verification Successful")
		BlockReceived = MsgReceived.BlockToSend
	}

	spew.Dump(BlockReceived)

	// Check whether Block is valid
	BlockValidity := isBlockValid(BlockReceived, Blockchain[len(Blockchain)-1])
	if BlockValidity {
		// Append to the main blockchain
		newBlockchain := append(Blockchain, BlockReceived)
		replaceChain(newBlockchain)
		writejson()
		for i := 0; i < len(profiles); i++ {
			if profiles[i].Role == 3 {
				// Append to the validator's blockchain
				BlockValidator(profiles[i].Userid, BlockReceived, i)
			}
		}
	}
	// Generate a Balance Block when the len(Blockchain)==4,9,14,19,24,...
	if (len(Blockchain)+1)%5 == 0 {
		GenerateBalanceBlock()
	}
}

func BlockValidator(uid int, BlockReceived Block, index int) {
	BlockValidity := isBlockValid(BlockReceived, profiles[index].Blockchain[len(profiles[index].Blockchain)-1])
	if BlockValidity {
		newBlockchain := append(profiles[index].Blockchain, BlockReceived)
		profiles[index].Blockchain = newBlockchain
		profiles[index].BlockchainHash = calculateBlockchainHash(profiles[index].Blockchain)
		writeUserjson()
	}
}

// Initiate a tcp port 3000 server, this function runs infinitely, it only needs to be run once every deployment
func StartLeaderServer() {
	nodeAddress1 := ":3000"
	ln1, err1 := net.Listen(protocol, nodeAddress1)
	if err1 != nil {
		log.Panic(err1)
	}
	defer ln1.Close()

	for {
		conn1, err1 := ln1.Accept()
		if err1 != nil {
			log.Panic(err1)
		}
		go ReadBlock(conn1)

	}
}
