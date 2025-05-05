// this file contains all the functionality related to database.
// adding data, get saved data

package main

import (
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"strconv"

	_ "github.com/microsoft/go-mssqldb"
)

var userid, password, server, database *string

func initFlags() {
	user := envData.msSQLUser
	pass := envData.msSQLPass
	svr := envData.msSQLServer
	dbname := envData.msSQLDBName

	userid = flag.String("U", user, "login_id")
	password = flag.String("P", pass, "password")
	server = flag.String("S", svr, "server_name[\\instance_name]")
	database = flag.String("d", dbname, "db_name")

}

func getDBConnection() (*sql.DB, error) {
	// if flags are already initialized, no need to initialize them again
	if flag.Lookup("U") == nil {
		initFlags()
	}
	flag.Parse()

	dsn := "server=" + *server + ";user id=" + *userid + ";password=" + *password + ";database=" + *database
	db, err := sql.Open("mssql", dsn)
	if err != nil {
		//fmt.Println("Cannot connect DB: Error Message: " + err.Error())
		log.Println("Cannot connect DB: Error Message: " + err.Error())
		return db, err
	} /*else {
		//fmt.Println("DB Connected.")
		log.Println("DB Connected.")
	}*/

	err = db.Ping()
	if err != nil {
		db.Close()
		//fmt.Println("Unable to ping DB: Error Message: " + err.Error())
		log.Println("Unable to ping DB: Error Message: " + err.Error())
		return db, err
	} else {
		//fmt.Println("Ping successful.")
		log.Println("Ping successful.")
	}

	return db, nil
}

func checkDBConnection() error {
	db, err := getDBConnection()
	if err != nil {
		return err
	}

	db.Close()

	return nil
}

func writeUserInfoIntoDB(id int, name string) (string, error) {
	retString := ""
	db, err := getDBConnection()
	if err != nil {
		//fmt.Println("Error Message: ", err.Error())
		log.Println("Error Message: ", err.Error())

		// structure for user data
		type Response struct {
			RowsInserted int    `json:"RowsInserted"`
			Msg          string `json:"Message"`
		}
		var response Response
		response.Msg = "There is some error. Please check logs for more details."

		jsonBytes, err := json.MarshalIndent(response, "", " ")
		if err != nil {
			return retString, err
		}
		retString = string(jsonBytes)

		return retString, err
	}

	defer db.Close()

	// insert data into DB table
	query := "INSERT INTO " + envData.msSQLDBName + ".dbo.UserInfo(ID, NAME) VALUES(?,?)"
	retString, err = execute(db, query, id, name)

	return retString, err
}

func readDatafromDB() (string, error) {
	retString := ""
	db, err := getDBConnection()
	if err != nil {
		//fmt.Println("Error Message: ", err.Error())
		log.Println("Error Message: ", err.Error())

		// structure for user data
		type UserInfo struct {
			ID   int    `json:"ID"`
			NAME string `json:"NAME"`
			Msg  string `json:"Message"`
		}

		var data UserInfo
		data.Msg = "There is some error. Please check logs for more details."

		jsonBytes, err := json.MarshalIndent(data, "", " ")
		if err != nil {
			return retString, err
		}
		retString = string(jsonBytes)

		return retString, nil
	}

	defer db.Close()

	query := "SELECT ID, NAME FROM " + envData.msSQLDBName + ".dbo.UserInfo"
	retString, err = queryDB(db, query, true, 0)
	if err != nil {
		return retString, err
	}

	return retString, nil
}

func readUserInfo(id int) (string, error) {
	retString := ""
	db, err := getDBConnection()
	if err != nil {
		//fmt.Println("Error Message: ", err.Error())
		log.Println("Error Message: ", err.Error())

		// structure for user data
		type UserInfo struct {
			ID   int    `json:"ID"`
			NAME string `json:"NAME"`
			Msg  string `json:"Message"`
		}

		var data UserInfo
		data.Msg = "There is some error. Please check logs for more details."

		jsonBytes, err := json.MarshalIndent(data, "", " ")
		if err != nil {
			return retString, err
		}
		retString = string(jsonBytes)

		return retString, nil
	}

	defer db.Close()

	query := "SELECT ID, NAME FROM " + envData.msSQLDBName + ".dbo.UserInfo WHERE ID=?"
	retString, err = queryDB(db, query, false, id)
	if err != nil {
		return retString, err
	}

	return retString, nil
}

func execute(db *sql.DB, query string, id int, name string) (string, error) {
	// structure for user data
	type Response struct {
		RowsInserted int    `json:"RowsInserted"`
		Msg          string `json:"Message"`
	}

	var response Response

	result, err := db.Exec(query, id, name)
	if err != nil {
		//fmt.Println("Error Message: ", err.Error())
		log.Println("Error Message: ", err.Error())
		response.RowsInserted = 0
		response.Msg = "There is some error. Please check logs for more details."
	} else {
		rowcount, _ := result.RowsAffected()
		//fmt.Println("No of rows afftected: ", rowcount)
		log.Println("No of rows afftected: ", rowcount)

		response.RowsInserted = int(rowcount)
		response.Msg = "User successfully added."
	}

	retString := ""
	jsonBytes, err := json.MarshalIndent(response, "", " ")
	if err != nil {
		return retString, err
	}
	retString = string(jsonBytes)

	return retString, nil
}

func queryDB(db *sql.DB, query string, bCompleteData bool, id int) (string, error) {
	retString := ""

	var rows *sql.Rows
	var err error

	if bCompleteData {
		rows, err = db.Query(query)
		if err != nil {
			return retString, err
		}
	} else {
		rows, err = db.Query(query, id)
		if err != nil {
			return retString, err
		}
	}

	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		return retString, err
	}

	if cols == nil {
		return retString, nil
	}

	vals := make([]interface{}, len(cols))
	for i := 0; i < len(cols); i++ {
		vals[i] = new(interface{})
	}

	// structure for user data
	type UserInfo struct {
		ID   int    `json:"ID"`
		NAME string `json:"NAME"`
		Msg  string `json:"Message"`
	}
	// slice of UserInfo structure
	dbData := make([]UserInfo, 0)

	count := 0
	for rows.Next() {
		count++

		err = rows.Scan(vals...)
		if err != nil {
			//fmt.Println(err)
			log.Println(err.Error())
			continue
		}

		var data UserInfo
		for i := 0; i < len(vals); i++ {
			strData := getRowValue(vals[i].(*interface{}))

			// var data UserInfo
			switch i {
			case 0: // ID
				data.ID, err = strconv.Atoi(strData)
				if err != nil {
					data.Msg = err.Error()
				}
			case 1: // DATA
				data.NAME = strData
			}
		}
		dbData = append(dbData, data)
	}

	if count == 0 {
		var data UserInfo
		if bCompleteData {
			//fmt.Println("No Data exist in database.")
			log.Println("No Data exist in database.")
			data.Msg = "No Data exist in database."
		} else {
			//fmt.Println("No Data exist in database for given ID.")
			log.Println("No Data exist in database for given ID: ", id)
			data.ID = id
			data.Msg = "No Data exist in database for given ID."
		}
		dbData = append(dbData, data)
	}

	jsonBytes, err := json.MarshalIndent(dbData, "", " ")
	if err != nil {
		return retString, err
	}
	retString = string(jsonBytes)

	return retString, nil
}

func getRowValue(pval *interface{}) string {
	retString := ""

	switch v := (*pval).(type) {
	case nil:
		//fmt.Print("NULL")
		log.Println("NULL")
	case int:
		//fmt.Print(v)
		log.Println(v)
	default:
		retString = fmt.Sprint(v)
	}
	return retString
}
