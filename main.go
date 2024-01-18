package main

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
	"shortly/db"
)

func main() {
	// Load the .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	conn, err := db.GetDBConnection()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	// ----- signup
	//user, err := models.SignUp(conn, "test2@test.io", "dsadsadsa")
	//if err != nil {
	//	fmt.Fprintf(os.Stderr, "Unable to signup user: %v\n", err)
	//	os.Exit(1)
	//}
	//fmt.Fprintf(os.Stderr, "Successfuly signup the user: %v\n", user.Id)

	// ----- signin by password
	//user, token, err := models.SignInByPassword(conn, "test@test.io", "asdasdasd")
	//if err != nil {
	//	fmt.Fprintf(os.Stderr, "Unable to signin user: %v\n", err)
	//	os.Exit(1)
	//}
	//fmt.Fprintf(os.Stderr, "Successfuly signin the user: %v, %v\n", user.Id, token.Value)

	// ----- signin by token
	//user, err := models.SignInByToken(conn, "6df1f6ea-0ac5-4cea-a3a4-ef50858720bb")
	//if err != nil {
	//	fmt.Fprintf(os.Stderr, "Unable to signin user: %v\n", err)
	//	os.Exit(1)
	//}
	//fmt.Fprintf(os.Stderr, "Successfuly signin the user: %v\n", user.Id)

	// ----- add url
	//url, err2 := models.AddURL(conn, user, "ya.ru")
	//if err2 != nil {
	//	fmt.Fprintf(os.Stderr, "Unable to create short URL: %v\n", err)
	//	os.Exit(1)
	//}
	//fmt.Fprintf(os.Stderr, "Successfuly create short URL: %v\n", url.Id)

	// ----- delete url
	//err3 := models.DeleteURL(conn, user, url.Id)
	//if err3 != nil {
	//	fmt.Fprintf(os.Stderr, "Unable to detele short URL: %v\n", err3)
	//	os.Exit(1)
	//}
	//fmt.Fprintf(os.Stderr, "Successfuly delete short URL: %v\n", url.Id)

	// ----- get url
	//url2, err3 := models.GetURL(conn, url.ShortLink)
	//if err3 != nil {
	//	fmt.Fprintf(os.Stderr, "Unable to get URL: %v\n", err3)
	//	os.Exit(1)
	//}
	//fmt.Fprintf(os.Stderr, "Successfuly get URL: %v\n", url2.Id)

	// ----- renew token
	//token2, err2 := models.RenewToken(conn, token.Value, user.Id)
	//if err2 != nil {
	//	fmt.Fprintf(os.Stderr, "Unable to renew token: %v\n", err2)
	//	os.Exit(1)
	//}
	//fmt.Fprintf(os.Stderr, "Successfuly renew token: %v, %v\n", user.Id, token2.Value)
}
