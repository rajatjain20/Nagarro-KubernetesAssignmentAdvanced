// This file contains all the handle functions
// for this backend application.

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

// structure for User Data
type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func root(w http.ResponseWriter, r *http.Request) {
	sEnv := envData.envName
	fmt.Fprintln(w, "Welcome to Kubernetes Advanced Assignment's backend application running on "+sEnv)
}

func checkDB(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*") // Add CORS header
	//w.WriteHeader(http.StatusOK)

	err := checkDBConnection()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Unable to connect to database.\nError Message:\n %s", err.Error())
		log.Println("checkDB(): Unable to connect to database.\nError Message:\n", err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Database is connected.")
	log.Println("checkDB(): Database is connected.")
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)

	fmt.Fprintf(w, "Server is live.")
	log.Println("healthCheck(): Server is live. response: ", http.StatusOK)
}

// to add user, this function is called from frontend
// It connencts to database based on the configuration (specific to environments)
func addUser(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		w.Header().Set("Access-Control-Allow-Origin", "*") // Add CORS header

		var user User
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			log.Println("addUser(): Error: ", err.Error())
			return
		}

		strData, err := writeUserInfoIntoDB(user.ID, user.Name)
		if err != nil {
			fmt.Fprintf(w, "Unable to add user into database.\nError Message: %s", err.Error())
			log.Println("addUser(): Unable to add user into database.\nError Message: ", err.Error())
		} else {
			fmt.Fprintln(w, strData)
			log.Println(strData)
		}

		//fmt.Fprintf(w, "User added: ID = %s, Name = %s", user.ID, user.Name)
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		log.Println("addUser(): Invalid request method")
	}
}

// to get user data, this function is called from frontend
// It connencts to database based on the configuration (specific to environments)
// and fetches data from database for provided specific user id and all data in one go.
func getUserInfo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*") // Add CORS header
	w.WriteHeader(http.StatusOK)

	if r.Method == http.MethodGet {
		val := r.URL.Query()
		if val.Has("id") {
			id, err := strconv.Atoi(val.Get("id"))
			if err != nil {

			}

			strData, err := readUserInfo(id)
			if err != nil {
				fmt.Fprintf(w, "Unable to retrive UserInfo.\nError Message: %s", err.Error())
				log.Println("getUserInfo(): Unable to retrive UserInfo.\nError Message: ", err.Error())
			} else {
				fmt.Fprintln(w, strData)
				log.Println(strData)
			}
		} else {
			strData, err := readDatafromDB()
			if err != nil {
				fmt.Fprintf(w, "Unable to retrive UserInfo.\nError Message: %s", err.Error())
				log.Println("getUserInfo(): Unable to retrive UserInfo.\nError Message: ", err.Error())
			}

			fmt.Fprintln(w, strData)
			log.Println(strData)
		}
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		log.Println("getUserInfo(): Invalid request method")
	}
}
