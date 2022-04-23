package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

/**
This file has all the code to:
- create a local user account file --> createUserAccountCopy
- write the data to the local user account file  -->  writeToUserLocalAcc
- read from the local user account  --> readLocalUserAccs
**/

const ACC_BALANCE_FILENAME = "userAccounts.json"

// create a locally stored blockchain
func createUserAccountCopy(dirname string) string {

	filename := dirname + "/" + ACC_BALANCE_FILENAME
	var _, err = os.Stat(filename)

	// create the local blockchain if it does not already exist
	if os.IsNotExist(err) {
		file, err := os.Create(filename)
		if err != nil {
			fmt.Println(err)

		}
		fmt.Println("User account file created successfully", filename)
		defer file.Close()
	} else {
		fmt.Println("User account file already exists!", filename)

	}

	return filename

}

// write to the blockchain in the file
// we always write it so mo need to update. writing is updating itself
func writeToUserLocalAcc(data []AccountBalance, fileDir string) {
	file, _ := json.MarshalIndent(data, "", " ")

	_ = ioutil.WriteFile(fileDir, file, 0644)
	fmt.Println("Local userAccount copy is now up to date")
}

// read the blocks from the local blockchain
func readLocalUserAccs(filepath string) []AccountBalance {
	file, _ := ioutil.ReadFile(filepath)

	localAccs := []AccountBalance{}

	_ = json.Unmarshal([]byte(file), &localAccs)

	//fmt.Println("Local user account balances as read", localAccs)
	return localAccs

}
