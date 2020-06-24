package main

import (
	d "first/app/models"
	routes "first/app/routes"
	util "first/app/utils"
	jobs "first/app/utils/cronJobs"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

// var psqlInfo = fmt.Sprintf("host=%s port=%s user=%s "+
// 	"dbname=%s sslmode=disable", GetEnvVariable("PGHOST"), GetEnvVariable("PGPORT"), GetEnvVariable("PGUSER"), GetEnvVariable("PGDBNAME"))

// func get(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(http.StatusOK)
// 	w.Write([]byte(`{"message": "get called"}`))
// }

// func post(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(http.StatusCreated)
// 	w.Write([]byte(`{"message": "post called"}`))
// }

// func put(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(http.StatusAccepted)
// 	w.Write([]byte(`{"message": "put called"}`))
// }

// func delete(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(http.StatusOK)
// 	w.Write([]byte(`{"message": "delete called"}`))
// }

func params(w http.ResponseWriter, r *http.Request) {
	pathParams := mux.Vars(r)
	w.Header().Set("Content-Type", "application/json")
	fmt.Printf("params-> %s", pathParams["userID"])
	userID := -1
	var err error
	if val, ok := pathParams["userID"]; ok {
		userID, err = strconv.Atoi(val)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"message": "need a number"}`))
			return
		}
	}

	query := r.URL.Query()
	w.Write([]byte(fmt.Sprintf(`{"userID": %d, "query": %s}`, userID, query)))
}

func getBatteryLevelInfo() {
	out, err := exec.Command("pmset", "-g", "batt").Output()
	if err != nil {
		log.Fatal(err)
	}
	type BatteryLevel struct {
		PowerType     string
		Level         string
		Status        string
		TimeRemaining string
	}
	// Get byte slice from string.
	bytes := []byte(out)
	s := string(bytes)
	val := strings.Split(s, ";")
	curVal := strings.Join(strings.Fields(val[0]), " ")
	curValLength := len(curVal)
	idxParen := strings.Index(curVal, ")")
	batterylevel := strings.Join(strings.Fields(curVal[idxParen+1:curValLength-1]), " ")
	thresholdVal := util.GetEnvVariable("BATTERY_LEVEL_THRESHOLD")
	threshold, err := strconv.Atoi(thresholdVal)
	chargingStatus := strings.Join(strings.Fields(val[1]), " ")

	if batt, err := strconv.Atoi(batterylevel); err == nil {
		if batt < threshold && chargingStatus == "discharging" {
			body := fmt.Sprintf("The current battery percentage is %d.", batt)
			util.SendEmail(body, "Your batttery level is running low.")
		}
	} else {
		log.Println(err)
	}

}

func main() {
	type User struct {
		Email string
		FName string
	}

	//var err error
	port := util.GetEnvVariable("PORT")
	curTime := time.Now()
	apiRouter := mux.NewRouter()
	api := apiRouter.PathPrefix("/api").Subrouter()
	// Connect to db
	db := d.InitPostgresDB()
	sqlStatement := `SELECT fname, email FROM users`
	var user User
	// Replace 3 with an ID from your database or another random
	// value to test the no rows use case.
	rows, err := db.Query(sqlStatement)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var myarray []string
	for rows.Next() {
		err := rows.Scan(&user.FName, &user.Email)
		if err != nil {
			log.Fatal(err)
		}
		myarray = append(myarray, fmt.Sprintf(`{"name": %s}`, user.Email))
	}
	//encjson, _ := json.Marshal(myarray)
	log.Println(myarray)
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
	//Run battery level cron
	jobs.ExecuteCronJob("@every 1m", getBatteryLevelInfo)

	// switch err := rows.Scan(&user.Email, &user.FName); err {
	// case sql.ErrNoRows:
	// 	fmt.Println("No rows were returned!")
	// 	defer db.Close()
	// case nil:
	// 	fmt.Println(user)
	// 	defer db.Close()
	// default:
	// 	panic(err)
	// }

	// fmt.Println("Successfully connected!")

	// api.HandleFunc("", get).Methods(http.MethodGet)
	// api.HandleFunc("", post).Methods(http.MethodPost)
	// api.HandleFunc("", put).Methods(http.MethodPut)
	// api.HandleFunc("", delete).Methods(http.MethodDelete)
	api.HandleFunc("/itunes/{name}", routes.GetItunesUserByName).Methods(http.MethodGet)
	api.HandleFunc("/user/{userID}", params).Methods(http.MethodGet)

	//for key, element := range os.Environ() {
	//	fmt.Println("Key:", key, "=>", "Element:", element)
	//}

	fmt.Printf("Go server listening on port %s at %s", port, curTime)

	log.Fatal(http.ListenAndServe(fmt.Sprintf("localhost:%s", port), apiRouter))

}
