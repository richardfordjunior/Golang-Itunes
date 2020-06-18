package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

//Load .env file
var itunesURL = GetEnvVariable("ITUNES_URL")

// GetItunesUserByName gets it artists songs and albums
func GetItunesUserByName(w http.ResponseWriter, r *http.Request) {

	pathParams := mux.Vars(r)

	response, err := http.Get(itunesURL + pathParams["name"])
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		http.Error(w, err.Error(), http.StatusBadRequest)
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf(`{"status":"OK","lastfetched":%s, "data": %s}`, time.Now(), data)))

	}
}
