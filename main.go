package main

//go:generate sqlboiler --wipe mysql

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
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
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
	domains, err := models.Admins().All(context.Background(), db)
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
	db.Exec("DELETE FROM mailbox")
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
	adminFoo := &models.Admin{Username: null.StringFrom("bar@foo.net"), Password: examplePassword, Super: true}
	err = adminFoo.Insert(context.Background(), db, boil.Infer())
	adminTest := &models.Admin{Username: null.StringFrom("info@test.de"), Password: examplePassword, Super: false}
	err = adminTest.Insert(context.Background(), db, boil.Infer())

	// Admin Domain Mapping
	fmt.Printf("******* Add Admin-Domain-Mappings *******\n")
	// TODO: models.DomainsAdmins seems to be missing
	err = adminFoo.AddDomains(context.Background(), db, false, domainFoo)
	dieIf(err)

	// Add a mailbox and link it to a domain
	fmt.Printf("******* Add Mailbox *******\n")
	fooMailbox := &models.Mailbox{
		Username:          "test@foo.net",
		Password:          examplePassword,
		Name:              null.StringFrom("Test User"),
		AltEmail:          null.StringFrom("test@test.com"),
		Quota:             0,
		LocalPart:         "test",
		Active:            false,
		AccessRestriction: "",
		Homedir:           null.StringFrom("/home/test"),
		Maildir:           null.String{},
		UID:               null.Int64From(5000),
		Gid:               null.Int64From(5000),
		HomedirSize:       null.Int64From(0),
		MaildirSize:       null.Int64From(0),
		SizeAt:            null.Time{},
		DeletePending:     null.BoolFrom(false),
		DomainID:          null.Int64From(domainFoo.ID),
	}
	err = fooMailbox.Insert(context.Background(), db, boil.Infer())
	dieIf(err)

	fmt.Printf("***** Seeding database with sample data, done. *****\n\n")

	// Now get a mailboxes and join them with the domain table
	var mods []qm.QueryMod
	mods = append(mods, models.MailboxWhere.ID.EQ(fooMailbox.ID))
	mods = append(mods, qm.Load(models.TableNames.Domain))
	log.Println("Get count of mailbox table")
	count, err := models.Mailboxes(mods...).Count(context.Background(), db)
	dieIf(err)

	mods = append(mods, qm.OrderBy("id"))
	log.Println("Get all mailboxes")
	mailboxes, err := models.Mailboxes(mods...).All(context.Background(), db)
	dieIf(err)
	log.Printf("Found %d mailboxes\n", count)
	for _, mailbox := range mailboxes {
		log.Printf("Found mailbox: %+v with domain %+v\n", mailbox, mailbox.R.Domain)
	}
}

func dieIf(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
