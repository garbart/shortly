package models

import (
	"context"
	"github.com/jackc/pgx/v5"
	"time"

	"github.com/google/uuid"
)

type Token struct {
	Id        int
	UserId    int
	Value     string
	ExpiredAt time.Time
}

func buildToken() Token {
	return Token{Value: uuid.New().String(), ExpiredAt: time.Now().AddDate(0, 0, 21)}
}

func RenewToken(conn *pgx.Conn, token string, userId int) (*Token, error) {
	// delete old token
	rows, err1 := conn.Query(context.Background(), "DELETE FROM shortly.tokens WHERE value = $1", token)
	if err1 != nil {
		return nil, err1
	}
	rows.Close()

	// create new token
	out := buildToken()
	out.UserId = userId
	err2 := conn.QueryRow(context.Background(), "INSERT INTO shortly.tokens (userid, value, expiredat) VALUES ($1, $2, $3) RETURNING id", out.UserId, out.Value, out.ExpiredAt).Scan(&out.Id)
	if err2 != nil {
		return nil, err2
	}

	return &out, nil
}
