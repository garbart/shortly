package models

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type Token struct {
	Id        int
	UserId    int
	Value     string
	ExpiredAt time.Time
}

func (token *Token) ToJSON() ([]byte, error) {
	return token.toJSON()
}

func (token *Token) toJSON() ([]byte, error) {
	resp := make(map[string]string)
	resp["id"] = strconv.Itoa(token.Id)
	resp["user_id"] = strconv.Itoa(token.UserId)
	resp["token"] = token.Value
	resp["expired_at"] = token.ExpiredAt.String()
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		return nil, err
	}
	return jsonResp, err
}

func buildToken() Token {
	return Token{Value: uuid.New().String(), ExpiredAt: time.Now().AddDate(0, 0, 21)}
}

func RenewToken(conn *pgx.Conn, token string) (*Token, error) {
	// get userId
	var userId int
	err0 := conn.QueryRow(context.Background(), "SELECT userId FROM shortly.tokens WHERE value = $1", token).Scan(&userId)
	if err0 != nil {
		return nil, err0
	}

	// delete old token
	rows, err1 := conn.Query(context.Background(), "DELETE FROM shortly.tokens WHERE value = $1 AND userId = $2", token, userId)
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
