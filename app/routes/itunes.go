package routes

import (
	util "first/app/utils"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

var itunesURL = util.GetEnvVariable("ITUNES_URL")

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
		data, error := ioutil.ReadAll(response.Body)
		if error == nil {
			util.APIResponse(w, data)
		} else {
			log.Fatal("Error reading response")
		}
	}
}
