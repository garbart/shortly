package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"shortly/api"
)

func main() {
	// Load the .env file
	envErr := godotenv.Load()
	if envErr != nil {
		log.Fatal("Error loading .env file")
	}

	// Setup HTTP handlers
	router := mux.NewRouter()
	//router.Handle("/api/signup", http.HandlerFunc(api.SignUpHandler))
	//router.Handle("/api/signin", http.HandlerFunc(api.SignInHandler))
	//router.Handle("/api/renew_token", http.HandlerFunc(api.RenewTokenHandler))
	//router.Handle("/api/url", http.HandlerFunc(api.URLHandler))
	router.Handle("/{short_link}", http.HandlerFunc(api.GetURLHandler))
	serverErr := http.ListenAndServe(":80", router)
	if serverErr != nil {
		fmt.Fprintf(os.Stderr, "Unable to setup server: %v\n", serverErr)
		os.Exit(1)
	}
}
