package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

/** This file has all the code to :
- store the blockchain locally
- update the locally stored chain
- read from the locally stored chain
**/

// create a locally stored blockchain
func createLocalBlockchainFile(dirname string) string {

	filename := dirname + "/blockchain.json"
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
