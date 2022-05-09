package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

// add household data to the database
func createEnergyData(w http.ResponseWriter, r *http.Request) {
	var fromFrontend EnergyPriceData
	var response EnergyDataResponse
	// get the data from json body
	decoder := json.NewDecoder(r.Body)
	// place the user data inside newUser
	if err := decoder.Decode(&fromFrontend); err != nil {
		fmt.Println("Failed adding a new user", err)
		respondWithJSON(w, r, http.StatusBadRequest, r.Body)
		return
	}
	defer r.Body.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var energyData EnergyPriceData = generateHouseholdEnergyPriceData()
	//write user info to the users collection
	if checkRepeatedDataEntry(fromFrontend.DateStr) {
		writeData, err := db.EnergyPriceHouse.InsertOne(ctx, energyData)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("New data added", writeData.InsertedID)
		response.Data = energyData
		respondWithJSON(w, r, http.StatusCreated, response)

	} else {
		fmt.Println("Fetching data from db")
		response.Data = getHouseholdData(fromFrontend.DateStr)
		respondWithJSON(w, r, http.StatusCreated, response)

	}

}

func generateHouseholdEnergyPriceData() EnergyPriceData {
	var energyData EnergyPriceData
	var maxMinDiff float64 = 0.02 // fluctuate price between +20% and -20%
	var elPriceHouse float64 = 0.221
	var maxElPriceHouse float64 = (elPriceHouse + (maxMinDiff * elPriceHouse))
	var minElPriceHouse float64 = (elPriceHouse - (maxMinDiff * elPriceHouse))
	//var elPriceBus float64 = 0.388

	//var maxPrice float64 =
	// Seed function to initialize the default Source since different behavior is required for each run.
	rand.Seed(time.Now().UnixNano())
	// Uniform random float (min <= x <= max)
	var data [24]float64
	var avg float64
	var summation float64 = 0
	for count := 0; count < 24; count++ {
		var elPrice float64 = (rand.Float64() * (maxElPriceHouse - minElPriceHouse)) + minElPriceHouse
		var finalPrice float64 = math.Floor(elPrice*1000) / 1000
		summation += finalPrice
		data[count] = finalPrice
		//fmt.Println("Price: ", finalPrice)
	}
	avg = summation / 23

	dateStr, weekday, nanoSec := getDateString()
	energyData.Day = weekday
	energyData.Average = avg
	energyData.Data = data
	energyData.DateStr = dateStr
	energyData.DateTime = nanoSec
	return energyData
}

func getDateString() (string, string, int64) {
	year, month, day := time.Now().Date()
	var monthStr string
	var dayStr string
	if month < 10 {
		monthStr = "0" + strconv.Itoa(int(month))
	} else {
		monthStr = strconv.Itoa(int(month))
	}

	if day < 10 {
		dayStr = "0" + strconv.Itoa(int(day))
	} else {
		dayStr = strconv.Itoa(int(day))
	}

	var dateStr string = dayStr + "-" + monthStr + "-" + strconv.Itoa(year)
	var weekday string = time.Now().Weekday().String()
	time := time.Now().Unix()
	//fmt.Println("Date String:", dateStr)
	return dateStr, weekday, time
}

func checkRepeatedDataEntry(dateStr string) bool {
	var result bool = false
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	// find username and password
	cursor, err := db.EnergyPriceHouse.Find(ctx, bson.M{"datestr": dateStr})
	if err != nil {
		log.Fatal(err)
	}

	var Profiles []EnergyPriceData
	if err = cursor.All(ctx, &Profiles); err != nil {
		log.Fatal(err)

	}
	if len(Profiles) == 0 {
		fmt.Println("No data on this date")
		result = true
	} else {
		fmt.Println("Data on this date exists")
		result = false
	}
	//fmt.Println(Profiles)
	return result
}

func getHouseholdData(dateStr string) EnergyPriceData {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	// find the relevant data
	cursor, err := db.EnergyPriceHouse.Find(ctx, bson.M{"datestr": dateStr})
	if err != nil {
		log.Fatal(err)
	}

	var Profiles []EnergyPriceData
	if err = cursor.All(ctx, &Profiles); err != nil {
		log.Fatal(err)

	}
	if len(Profiles) == 1 {
		return Profiles[0]

	} else {
		fmt.Println("No Data on this date exists")
		return Profiles[0]
	}
	//fmt.Println(Profiles)

}
