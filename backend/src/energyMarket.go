package main

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"time"
)

func createEnergyData() {
	var energyData EnergyPriceData = generateHouseholdEnergyPriceData()
	//write user info to the users collection
	writeUser, err := db.EnergyPriceHouse.InsertOne(mongoparams.ctx, energyData)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("New data added", writeUser.InsertedID)
}

func generateHouseholdEnergyPriceData() EnergyPriceData {
	var energyData EnergyPriceData
	var maxMinDiff float64 = 0.2 // fluctuate price between +20% and -20%
	var elPriceHouse float64 = 0.221
	var maxElPriceHouse float64 = (elPriceHouse + (maxMinDiff * elPriceHouse))
	var minElPriceHouse float64 = (elPriceHouse + (maxMinDiff * elPriceHouse))
	//var elPriceBus float64 = 0.388

	//var maxPrice float64 =
	// Seed function to initialize the default Source since different behavior is required for each run.
	rand.Seed(time.Now().UnixNano())
	// Uniform random float (min <= x <= max)
	var data [24]float64
	var avg float64
	var summation float64 = 0
	for count := 0; count < 24; count++ {
		var elPrice float64 = (rand.Float64() * (maxElPriceHouse - minElPriceHouse + 0.1)) + minElPriceHouse
		var finalPrice float64 = math.Floor(elPrice*1000) / 1000
		summation += finalPrice
		data[count] = finalPrice
		fmt.Println("Price: ", finalPrice)
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

// func processTime() {
// 	dateStr, weekday, nanoSec := getDateString()
// 	fmt.Println("Date String:", dateStr)
// 	fmt.Println("Week day:", weekday)
// 	fmt.Println("Time in nano seconds:", nanoSec)
// }

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
