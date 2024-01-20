package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"shortly/db"
	"shortly/models"
)

type RenewTokenRequest struct {
	Token string `json:"token"`
}

func RenewTokenHandler(w http.ResponseWriter, r *http.Request) {
	// check method
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// get DB connection
	conn, err := db.GetDBConnection()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer conn.Close(context.Background())

	// get token from request
	var request RenewTokenRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err2 := decoder.Decode(&request)
	if err2 != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// renew token
	token, err3 := models.RenewToken(conn, request.Token)
	if err3 != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	resp, err := token.ToJSON()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(resp)
}
