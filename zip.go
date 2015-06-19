package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
	"strconv"
)

var key string

type distanceResponse struct {
	Distance float64
}

type zipRenderContext struct {
	distanceStrings []string
}

func loadKey() string {
	data, err := ioutil.ReadFile("zip_api.json")
	if err != nil {
		fmt.Println(err)
	}
	err = json.Unmarshal(data, &key)
	if err != nil {
		fmt.Println(err)
	}
	return key
}

// %s/distance.json/%s/%s
func getDistance(fromZip string, toZip string) float64 {
	requestURL := "https://www.zipcodeapi.com/rest/" + key + "/distance.json/" + fromZip + "/" + toZip + "/miles"
	resp, err := http.Get(requestURL)
	if err != nil {
		log.Fatal(err) // TODO: COME BACK AND ADD ACTUAL ERROR HANDLING
		panic(err)
	}
	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)
	var distance distanceResponse
	err = decoder.Decode(&distance)
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Println(distance.Distance)
	return distance.Distance
}

func init() {
	key = loadKey()
}

func zipHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "zip")
}

func zipSubmitHandler(w http.ResponseWriter, r *http.Request) {
	fromZip := r.FormValue("fromZip")
	toZip := r.FormValue("toZip")
	fmt.Println(fromZip)
	fmt.Println(strings.Fields(toZip))
	for _, zip := range strings.Fields(toZip) {
		fmt.Printf(zip + " : ")
		distance := getDistance(fromZip, zip)
		fmt.Printf(strconv.FormatFloat(distance, 'f', -1, 64) + " miles.\n")
		time.Sleep(0.25 * time.Second)
	}
}
