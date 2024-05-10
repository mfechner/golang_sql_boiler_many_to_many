package main

//go:generate sqlboiler mysql

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/pgtype"
	_ "github.com/mattn/go-sqlite3"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"golang.org/x/crypto/bcrypt"

	"go_sql_boiler/internal/db/models"
)

func main() {
	var db *sql.DB
	boil.DebugMode = true

	cfg, err := LoadEnvVariables(".env")
	if err != nil {
		log.Fatal(err)
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=true&loc=Local", cfg.DBUsername, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	// add sample data
	InsertSampleData(db)

	// retrieve all domain and print them to console

	fmt.Printf("Run query to domain table\n")
	domains, err := models.Admins(Load(models.Admins.Domains)).All(context.Background(), db)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("done. %+v\n", domains)

	for _, domain := range domains {
		fmt.Printf("domain: %#v\n", domain)
	}

}

func InsertSampleData(db *sql.DB) {
	fmt.Printf("\n***** Start to seed the database with test data *****\n")
	examplePlainPassword := "test"
	bytesPassword, _ := bcrypt.GenerateFromPassword([]byte(examplePlainPassword), bcrypt.DefaultCost)
	examplePassword := string(bytesPassword)

	db.Exec("DELETE FROM domain_admins")
	db.Exec("DELETE FROM domain")
	db.Exec("DELETE FROM admin")
	fmt.Printf("******* Add Domains *******\n")
	domainTest := &models.Domain{ID: 3, Domain: "test.de"}
	err := domainTest.Insert(context.Background(), db, boil.Infer())
	dieIf(err)
	domainFoo := &models.Domain{ID: 4, Domain: "foo.net", Description: null.String{String: "Domain from Foo"}}
	err = domainFoo.Insert(context.Background(), db, boil.Infer())
	dieIf(err)

	// Create admins
	fmt.Printf("******* Add Admins *******\n")
	adminFoo := &models.Admin{Username: "bar@foo.net", Password: examplePassword, Super: true}
	err = adminFoo.Insert(context.Background(), db, boil.Infer())
	adminTest := &models.Admin{Username: "info@test.de", Password: examplePassword, Super: false}
	err = adminTest.Insert(context.Background(), db, boil.Infer())

	// Admin Domain Mapping
	fmt.Printf("******* Add Admin-Domain-Mappings *******\n")
	// TODO: models.DomainsAdmins seems to be missing

	err = domainTest.AddAdminAdmins(context.Background(), db, true, adminTest)

	fmt.Printf("***** Seeding database with sample data, done. *****\n\n")
}

func dieIf(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
