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

type SignUpRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func SignUpHandler(w http.ResponseWriter, r *http.Request) {
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

	// get data from request
	var request SignUpRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err2 := decoder.Decode(&request)
	if err2 != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// signup user
	_, err = models.SignUp(conn, request.Email, getMD5Hash(request.Password))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func SignInHandler(w http.ResponseWriter, r *http.Request) {
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

	// auth user
	_, token, err := basicAuth(conn, r)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Cannot auth user: %v\n", err)
		w.WriteHeader(http.StatusUnauthorized)
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
