package main

import (
	"crypto/rsa"
	"database/sql"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/rs/cors"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

// Secret to generate token
const Secret = "&nPKf2d(/u]/d5k8"

var (
	settings *mainSettings
	signKey  *rsa.PrivateKey
	db       *sql.DB
)

func init() {
	var err error
	settings, err = loadSettings()

	if err != nil {
		log.Fatal(err)
	}

	// Get private key
	signBytes, err := ioutil.ReadFile(settings.KeyPath)
	if err != nil {
		log.Fatal(err)
	}

	signKey, err = jwt.ParseRSAPrivateKeyFromPEM(signBytes)
	if err != nil {
		log.Fatal(err)
	}

	// Connect to DB
	db, err = sql.Open("postgres", settings.DatabaseConnect)

	if err != nil {
		log.Fatal(err)
	}

	db.SetMaxIdleConns(10)
}

func main() {
	defer db.Close()

	r := mux.NewRouter()
	r.HandleFunc("/login", login).Methods(http.MethodPost)
	r.HandleFunc("/logout", logout).Methods(http.MethodPost)

	handler := cors.Default().Handler(r)

	log.Println("Listening for connections on port: ", settings.Port)
	log.Fatal(http.ListenAndServe(":"+settings.Port, handler))
}
