package main

import (
	"log"
	"os"

	mysqlCfg "github.com/go-sql-driver/mysql"

	"github.com/andresdal/gobackexample/config"
	"github.com/andresdal/gobackexample/db"

	"github.com/golang-migrate/migrate/v4/database/mysql"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	db, err := db.NewMySQLStorage(mysqlCfg.Config{
		User: config.Envs.DBUser,
		Passwd: config.Envs.DBPassword,
		Net: "tcp",
		Addr: config.Envs.DBAddress,
		DBName: config.Envs.DBName,
		AllowNativePasswords: true,
		ParseTime: true,
	})
	if err != nil {
		log.Fatal(err)
	}

	driver, err := mysql.WithInstance(db, &mysql.Config{})
	if err != nil {
		log.Fatal(err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://cmd/migrate/migrations",
		"mysql",
		driver,
	)
	if err != nil {
		log.Fatal(err)
	}

	cmd := os.Args[(len(os.Args)-1)]
	if cmd == "up" {
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Fatal(err)
		}
	}
	if cmd == "down" {
		if err := m.Down(); err != nil && err != migrate.ErrNoChange {
			log.Fatal(err)
		}
	}
}