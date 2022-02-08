package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"time"
)

func runEnergyForecast(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Methods", "GET,POST,OPTIONS,DELETE,PUT")
	var userId string

	// get the data from the request body
	decoder := json.NewDecoder(r.Body)
	// place the user data inside newUser
	if err := decoder.Decode(&userId); err != nil {
		fmt.Println("Failed adding new forecast data", err)
		respondWithJSON(w, r, http.StatusBadRequest, r.Body)
		return
	}
	defer r.Body.Close()
	fmt.Println("This is the user id for forecasting from frontend", userId)
	// to prevent backend to timeout
	mongoparams.ctx, mongoparams.cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer mongoparams.cancel()

	cmd := exec.Command("python", "expSmoothingForecast.py", userId)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		fmt.Println("This shit does not work", err)
	} else {
		fmt.Println("Ran py script successfully")
	}
	var newRequest = ""
	respondWithJSON(w, r, http.StatusCreated, newRequest)
}

// func runEnergyForecast(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Add("Access-Control-Allow-Methods", "GET,POST,OPTIONS,DELETE,PUT")
// 	var newRequest = ""
// 	//c := exec.Command("python", "expSmoothingForecast.py")
// 	c := exec.Command("python expSmoothingForecast.py", "raf")
// 	err := c.Run()
// 	if err != nil {
// 		fmt.Println("Error: ", err)
// 	}
// 	//fmt.Println(out)
// 	// if err := c.Output(); err != nil {
// 	// 	fmt.Println("Error: ", err)

// 	// 	respondWithJSON(w, r, http.StatusCreated, "Shit doesn't work")
// 	// }
// 	fmt.Println("Ran py script successfully")
// 	respondWithJSON(w, r, http.StatusCreated, newRequest)
// }
