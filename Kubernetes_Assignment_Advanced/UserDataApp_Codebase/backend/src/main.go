//go:debug x509negativeserial=1
package main

import (
	"log"
	"net/http"
	"os"
	"runtime"

	"github.com/joho/godotenv"
)

type envVarData struct {
	envName, port                                  string
	msSQLServer, msSQLDBName, msSQLUser, msSQLPass string
}

var envData envVarData

func main() {
	// Load Environment variable data
	initData(&envData)

	// Handle Functions
	http.HandleFunc("/checkDB", checkDB)     //readiness
	http.HandleFunc("/healthz", healthCheck) //liveness
	http.HandleFunc("/addUser", func(w http.ResponseWriter, r *http.Request) {
		// Handle CORS preflight request
		if r.Method == http.MethodOptions {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
			return
		}
		addUser(w, r)
	})
	http.HandleFunc("/getUserInfo", getUserInfo)
	http.HandleFunc("/", root)

	//fmt.Println("Server is running at port: ", envData.port)
	log.Println("Welcome to Kubernetes Advanced Assignment's backend application running on " + envData.envName)
	log.Println("Server is running at port: ", envData.port)
	http.ListenAndServe(":"+envData.port, nil)
}

// check the OS
func isWindowsOS() bool {
	return runtime.GOOS == "windows"
}

// Data will be read from .env file on windows and from env variables on ubuntu
func initData(envData *envVarData) {
	if isWindowsOS() {
		envFile, _ := godotenv.Read("..\\config\\.env")
		envData.envName = envFile["ENV_NAME"]
		envData.port = envFile["PORT"]
		envData.msSQLServer = envFile["MSSQL_SERVER"]
		envData.msSQLDBName = envFile["MSSQL_DBNAME"]
		envData.msSQLUser = envFile["MSSQL_USER"]
		envData.msSQLPass = envFile["MSSQL_PASS"]
	} else {
		envData.envName = os.Getenv("ENV_NAME")
		envData.port = os.Getenv("PORT")
		envData.msSQLServer = os.Getenv("MSSQL_SERVER")
		envData.msSQLDBName = os.Getenv("MSSQL_DBNAME")
		envData.msSQLUser = os.Getenv("MSSQL_USER")
		envData.msSQLPass = os.Getenv("MSSQL_PASS")
	}
}
