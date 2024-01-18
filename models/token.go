package models

import "time"

type Token struct {
	Id        int
	UserId    int
	Value     string
	ExpiredAt time.Time
}

func RenewToken(token Token) (*Token, error) {
	return nil, nil
}
