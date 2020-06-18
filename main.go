package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
  "github.com/richardfordjunior/first"
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

func main() {
	type User struct {
		Email string
		FName string
	}

	//var err error
	port := GetEnvVariable("PORT")
	curTime := time.Now()
	apiRouter := mux.NewRouter()
	api := apiRouter.PathPrefix("/api").Subrouter()
	// Connect to db
	db := InitPostgresDB()
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
	api.HandleFunc("/itunes/{name}", GetItunesUserByName).Methods(http.MethodGet)
	api.HandleFunc("/user/{userID}", params).Methods(http.MethodGet)

	//for key, element := range os.Environ() {
	//	fmt.Println("Key:", key, "=>", "Element:", element)
	//}

	fmt.Printf("Go server listening on port %s at %s", port, curTime)

	log.Fatal(http.ListenAndServe(fmt.Sprintf("localhost:%s", port), apiRouter))

}
