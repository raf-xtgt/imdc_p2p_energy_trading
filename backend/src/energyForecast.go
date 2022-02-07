package main

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
)

func runEnergyForecast(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Methods", "GET,POST,OPTIONS,DELETE,PUT")
	var newRequest = ""

	cmd := exec.Command("python", "expSmoothingForecast.py", "raf")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		fmt.Println("This shit does not work", err)
	} else {
		fmt.Println("Ran py script successfully")
	}
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
