package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

func runBuyEnergyForecast(w http.ResponseWriter, r *http.Request) {
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

	// executing the python script from golang
	cmd := exec.Command("python", "buyForecastForAllUsers.py")
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

// function to get the latest buy energy forecast
func getLatestBuyForecast(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Methods", "GET,POST,OPTIONS,DELETE,PUT")
	var dateStr string

	// get the data from the request body
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&dateStr); err != nil {
		fmt.Println("Failed getting latest forecast data", err)
		respondWithJSON(w, r, http.StatusBadRequest, r.Body)
		return
	}
	defer r.Body.Close()
	fmt.Println("This is the date for getting buy forecast data from frontend", dateStr)
	// to prevent backend to timeout
	mongoparams.ctx, mongoparams.cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer mongoparams.cancel()

	// find data
	cursor, err := db.BuyOrderForecast.Find(mongoparams.ctx, bson.M{"date": dateStr})
	if err != nil {
		log.Fatal(err)
	}

	// put the data in forecast response
	var buyOrderForecastResponse []BuyForecastResponse
	if err = cursor.All(mongoparams.ctx, &buyOrderForecastResponse); err != nil {
		fmt.Println("Got error here")
		log.Fatal(err)
	}

	respondWithJSON(w, r, http.StatusCreated, buyOrderForecastResponse)
}
