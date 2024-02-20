package main

import (
	"bytes"
	"encoding/csv"
	"flag"
	"fmt"
	"inspiro_quotes_web/internal/quotes"
	"log"
	"math/rand"
	"os"
	"strings"

	"golang.org/x/term"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

func main() {
	var action string

	flag.StringVar(&action, "action", "init", `Action to perform. Options are:
	- init: Initialize the database
	- create-admin: Create an admin user`)
	flag.Parse()

	switch action {
	case "init":
		InitializeDB()
	case "create-admin":
		CreateAdmin()
	default:
		log.Fatal("Invalid action")
	}

}

func init() {
	var err error
	db, err = gorm.Open(sqlite.Open("quotes.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&quotes.Quote{}, &quotes.SuperUser{}, &quotes.ImageForQuote{})
}

func InitializeDB() {
	file, err := os.Open("quotes.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = ';'
	_, err = reader.Read()
	if err != nil {
		log.Fatal(err)
	}

	reader.FieldsPerRecord = 3
	records, err := reader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}
	quotes_list := make([]*quotes.Quote, len(records))
	for i, record := range records {
		quotes_list[i] = &quotes.Quote{
			Text:   record[0],
			Author: record[1],
			Genre:  strings.ToLower(record[2]),
			Lang:   "en",
		}
	}
	//Save the quotes to the database
	threeRandomQuotes := []quotes.Quote{*quotes_list[rand.Intn(len(quotes_list))], *quotes_list[rand.Intn(len(quotes_list))], *quotes_list[rand.Intn(len(quotes_list))]}
	num_errors := 0
	for _, q := range threeRandomQuotes {
		var quo quotes.Quote
		err := db.Where("text = ?", q.Text).First(&quo).Error
		if err == gorm.ErrRecordNotFound {
			num_errors++
		}
	}
	if num_errors == 3 {
		quotes.CreateManyQuotes(quotes_list)
	}
}

func CreateAdmin() {
	var superUser quotes.SuperUser
	var username string
	log.Println("Creating admin user")
	log.Println("Enter username: ")
	_, err := fmt.Scanln(&username)
	if err != nil {
		log.Fatal(err)
	}
	superUser.Username = username
	log.Println("Enter password: ")
	password1, err := term.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Re-enter password: ")
	password2, err := term.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		log.Fatal(err)
	}
	if len(password1) <= 8 || len(password2) <= 8 {
		log.Fatal("Inavlid password length")
	}
	if len(password1) != len(password2) || !bytes.Equal(password1, password2) {
		log.Fatal("Passwords do not match")
	}
	err = superUser.SetPasswordFromBytes(password1)
	if err != nil {
		log.Fatal(err)
	}
	db.Create(&superUser)
}
