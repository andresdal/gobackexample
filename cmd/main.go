package main

import (
	"database/sql"
	"log"

	"github.com/andresdal/gobackexample/cmd/api"
	"github.com/andresdal/gobackexample/config"
	"github.com/andresdal/gobackexample/db"
	"github.com/go-sql-driver/mysql"
)

func main() {
	db, err := db.NewMySQLStorage(mysql.Config{
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

	initStorage(db)

	server := api.NewAPIServer(":8080", db)
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}

func initStorage(db *sql.DB) {
	err := db.Ping()
	if err != nil {
		log.Fatal("Cannot connect to DB:", err)
	}

	log.Println("Connected to DB successfully")
}