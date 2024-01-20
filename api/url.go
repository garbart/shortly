package api

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"shortly/models"

	"shortly/db"

	"github.com/gorilla/mux"
)

type GetURLRequest struct {
	ShortLink string `json:"short_link"`
}

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

	// redirect
	http.Redirect(w, r, url.OriginalLink, http.StatusSeeOther)
}
