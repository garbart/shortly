package models

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"log"

	"github.com/jackc/pgx/v5"
)

type URL struct {
	Id           int
	UserId       int
	OriginalLink string
	ShortLink    string
	Views        int
}

func getUserUrls(conn *pgx.Conn, userId int) ([]URL, error) {
	rows, err := conn.Query(context.Background(), "SELECT id, userId, originalLink, shortLink, views FROM shortly.urls WHERE userid = $1", userId)
	if err != nil {
		return []URL{}, err
	}
	defer rows.Close()

	var urls []URL
	for rows.Next() {
		var url URL
		err := rows.Scan(&url.Id, &url.UserId, &url.OriginalLink, &url.ShortLink, &url.Views)
		if err != nil {
			log.Fatal(err)
		}
		urls = append(urls, url)
	}
	if err := rows.Err(); err != nil {
		return []URL{}, err
	}

	return urls, nil
}

func GetURL(conn *pgx.Conn, shortLink string) (*URL, error) {
	// select url
	var out URL
	err2 := conn.QueryRow(context.Background(), "SELECT id, userId, originalLink, shortLink, views  FROM shortly.urls WHERE shortlink = $1", shortLink).Scan(&out.Id, &out.UserId, &out.OriginalLink, &out.ShortLink, &out.Views)
	if err2 != nil {
		return nil, err2
	}

	return &out, nil
}

func AddURL(conn *pgx.Conn, user *User, originalLink string) (*URL, error) {
	// generate short link
	b := make([]byte, 4)
	_, err1 := rand.Read(b)
	if err1 != nil {
		return nil, err1
	}
	shortLink := hex.EncodeToString(b)

	// insert into db
	var id int
	err2 := conn.QueryRow(context.Background(), "INSERT INTO shortly.urls (userid, originallink, shortlink) VALUES ($1, $2, $3) RETURNING id", user.Id, originalLink, shortLink).Scan(&id)
	if err2 != nil {
		return nil, err2
	}

	// add url to user object
	out := URL{Id: id, UserId: user.Id, OriginalLink: originalLink, ShortLink: shortLink, Views: 0}
	user.Urls = append(user.Urls, out)

	return &out, nil
}

func DeleteURL(conn *pgx.Conn, user *User, urlId int) error {
	// delete url from db
	rows, err := conn.Query(context.Background(), "DELETE FROM shortly.urls WHERE id = $1 AND userId = $2", urlId, user.Id)
	if err != nil {
		return err
	}
	rows.Close()

	// delete url from user object
	for i, other := range user.Urls {
		if other.Id == urlId {
			user.Urls = append(user.Urls[:i], user.Urls[i+1:]...)
			break
		}
	}

	return nil
}
