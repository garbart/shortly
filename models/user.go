package models

import (
	"context"
	"github.com/google/uuid"
	"time"

	_ "github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type User struct {
	Id           int
	Email        string
	PasswordHash string
	Urls         []URL
}

func SignUp(conn *pgx.Conn, email, passwordHash string) (*User, error) {
	var id int
	err := conn.QueryRow(context.Background(), "INSERT INTO shortly.users (email, passwordhash) VALUES ($1, $2) RETURNING id", email, passwordHash).Scan(&id)
	if err != nil {
		return nil, err
	}

	out := User{Id: id, Email: email, PasswordHash: passwordHash}
	return &out, nil
}

func SignInByPassword(conn *pgx.Conn, email, passwordHash string) (*User, *Token, error) {
	// select user
	user := User{}
	err1 := conn.QueryRow(context.Background(), "SELECT id, email, passwordHash FROM shortly.users WHERE email = $1 AND passwordhash = $2", email, passwordHash).Scan(&user.Id, &user.Email, &user.PasswordHash)
	if err1 != nil {
		return nil, nil, err1
	}

	// select urls
	urls, err2 := getUserUrls(conn, user.Id)
	if err2 != nil {
		return nil, nil, err2
	}
	user.Urls = urls

	// insert token
	token := Token{Value: uuid.New().String(), ExpiredAt: time.Now().AddDate(0, 0, 21)}
	err3 := conn.QueryRow(context.Background(), "INSERT INTO shortly.tokens (userid, value, expiredat) VALUES ($1, $2, $3) RETURNING id", user.Id, token.Value, token.ExpiredAt).Scan(&token.Id)
	if err3 != nil {
		return nil, nil, err3
	}

	return &user, &token, nil
}

func SignInByToken(conn *pgx.Conn, tokenValue string) (*User, error) {
	// select token
	token := Token{}
	err := conn.QueryRow(context.Background(), "SELECT id, userid, value, expiredat FROM shortly.tokens WHERE value = $1 AND expiredat > $2", tokenValue, time.Now()).Scan(&token.Id, &token.UserId, &token.Value, &token.ExpiredAt)
	if err != nil {
		return nil, err
	}

	// select user
	user := User{}
	err1 := conn.QueryRow(context.Background(), "SELECT id, email, passwordHash FROM shortly.users WHERE id = $1", token.Id).Scan(&user.Id, &user.Email, &user.PasswordHash)
	if err1 != nil {
		return nil, err1
	}

	// select urls
	urls, err2 := getUserUrls(conn, user.Id)
	if err2 != nil {
		return nil, err2
	}
	user.Urls = urls

	return &user, nil
}
