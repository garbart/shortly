package api

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"net/http"
	"strings"

	"shortly/models"

	"github.com/jackc/pgx/v5"
)

func getMD5Hash(text string) string {
	hash := md5.Sum([]byte(text + "8F5UsjhYlPzMRCoQ69hd"))
	return hex.EncodeToString(hash[:])
}

func basicAuth(conn *pgx.Conn, r *http.Request) (*models.User, *models.Token, error) {
	// get 'Authorization' header
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return nil, nil, errors.New("no 'Authorization' header")
	}

	// get email and password
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || parts[0] != "Basic" {
		return nil, nil, errors.New("bad 'Authorization' format")
	}
	decoded, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, nil, errors.New("bad auth data")
	}
	credentials := strings.SplitN(string(decoded), ":", 2)
	if len(credentials) != 2 {
		return nil, nil, errors.New("bad auth data")
	}

	// get user
	user, token, err := models.SignInByPassword(conn, credentials[0], getMD5Hash(credentials[1]))
	if err != nil {
		return nil, nil, err
	}
	return user, token, nil
}

func bearerAuth(conn *pgx.Conn, r *http.Request) (*models.User, error) {
	// get 'Authorization' header
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return nil, errors.New("no 'Authorization' header")
	}

	// get token
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		return nil, errors.New("'Authorization' header bad format")
	}
	token := parts[1]

	// get user
	user, err := models.SignInByToken(conn, token)
	if err != nil {
		return nil, err
	}
	return user, nil
}
