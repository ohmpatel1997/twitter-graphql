package database

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/mysql"
	_ "github.com/golang-migrate/migrate/source/file"
	"log"
)

var Db *sql.DB

func InitDB() {
	db, err := sql.Open("mysql", "root:root12345@tcp(docker.for.mac.localhost)/twitter?parseTime=true")
	if err != nil {
		log.Println(err)
		return
	}

	if err = db.Ping(); err != nil {
		log.Println(err)
		return
	}
	Db = db
}

func Migrate() error {

	if err := Db.Ping(); err != nil {
		log.Println(err)
		return err
	}
	driver, _ := mysql.WithInstance(Db, &mysql.Config{})
	m, _ := migrate.NewWithDatabaseInstance(
		"file://internal/pkg/db/migrations/mysql",
		"mysql",
		driver,
	)
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Println(err)
		return err
	}
	return nil
}
