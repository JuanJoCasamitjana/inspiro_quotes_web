package main

import (
	"bytes"
	"encoding/csv"
	"flag"
	"fmt"
	"inspiro_quotes_web/internal/quotes"
	"inspiro_quotes_web/internal/server"
	"log"
	"math/rand"
	"os"
	"strings"

	"github.com/glebarez/sqlite"
	"golang.org/x/term"
	"gorm.io/gorm"
)

/* func main() {
	server.SetupAndRun()
} */

func main() {
	var action string
	var port string

	flag.StringVar(&action, "action", "run", `Action to perform. Options are:
	- init: Initialize the database
	- create-admin: Create an admin user
	- run: Run the server
	- init-run: Initialize the database and run the server
	- full: Initialize the database, create an admin user from env variables and run the server
Default is run.`)
	flag.StringVar(&port, "port", ":8080", "Port to run the server on")
	flag.Parse()

	env_action := os.Getenv("ACTION")
	if env_action != "" {
		action = env_action
	}

	switch action {
	case "init":
		InitializeDB()
	case "create-admin":
		CreateAdmin()
	case "run":
		server.SetupAndRun(port)
	case "init-run":
		InitializeDB()
		server.SetupAndRun(port)
	case "full":
		server.SetupAndRun(port)
		InitializeDB()
		CreateAdminFromEnv()
	default:
		log.Fatal("Invalid action")
	}

}

func InitializeDB() {
	var err error
	DB, err := gorm.Open(sqlite.Open("quotes.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	DB.AutoMigrate(&quotes.Quote{}, &quotes.SuperUser{}, &quotes.ImageForQuote{})
	log.Println("Initializing database")
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
	log.Println("Creating quotes, total:", len(records))
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
		err := DB.Where("text = ?", q.Text).First(&quo).Error
		if err == gorm.ErrRecordNotFound {
			num_errors++
		}
	}
	log.Println("Saving all quotes to the database")
	if num_errors == 3 {
		DB.Transaction(func(tx *gorm.DB) error {
			if err := tx.CreateInBatches(&quotes_list, 200).Error; err != nil {
				return err
			}
			return nil
		})
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
	quotes.DB.Create(&superUser)
}

func CreateAdminFromEnv() {
	var superUser quotes.SuperUser
	superUser.Username = os.Getenv("ADMIN_USERNAME")
	password := []byte(os.Getenv("ADMIN_PASSWORD"))
	if len(password) <= 8 || len(superUser.Username) == 0 {
		log.Fatal("Invalid password or username")
	}
	err := superUser.SetPasswordFromBytes(password)
	if err != nil {
		log.Fatal(err)
	}
	quotes.DB.Create(&superUser)
}
