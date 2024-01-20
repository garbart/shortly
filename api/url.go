package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"shortly/db"
	"shortly/models"

	"github.com/gorilla/mux"
)

func GetURLHandler(w http.ResponseWriter, r *http.Request) {
	// check method
	if r.Method != http.MethodGet {
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

	// get link from request
	vars := mux.Vars(r)
	shortUrl := vars["short_link"]
	if shortUrl == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// get original link from db
	url, err1 := models.GetURL(conn, shortUrl)
	if err1 != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// increment views
	_ = models.IncrementURLViews(conn, url.Id)

	// redirect
	http.Redirect(w, r, url.OriginalLink, http.StatusSeeOther)
}

func GetAllURLsHandler(w http.ResponseWriter, r *http.Request) {
	// check method
	if r.Method != http.MethodGet {
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
	user, err := bearerAuth(conn, r)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Cannot auth user: %v\n", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// get all user urls
	urls, err := models.GetUserUrls(conn, user.Id)

	resp := make(map[string]interface{})
	resp["urls"] = urls
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonResp)
}

type CreateURLRequest struct {
	OriginalLink string `json:"original_link"`
}

func CreateURLHandler(w http.ResponseWriter, r *http.Request) {
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
	user, err := bearerAuth(conn, r)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Cannot auth user: %v\n", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// get request body
	var request CreateURLRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&request)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// create short URL
	url, err := models.AddURL(conn, user, request.OriginalLink)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// setup response body
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	resp, err := url.ToJSON()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(resp)
}

func DeleteURLHandler(w http.ResponseWriter, r *http.Request) {
	// check method
	if r.Method != http.MethodDelete {
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
	user, err := bearerAuth(conn, r)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Cannot auth user: %v\n", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// get shortUrl from query
	shortLink := r.URL.Query().Get("short_link")
	if shortLink == "" {
		fmt.Fprintf(os.Stderr, "No 'short_link' param in query: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// delete url
	err = models.DeleteURL(conn, user, shortLink)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
