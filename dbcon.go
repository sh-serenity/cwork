package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)
var db *sql.DB
func dbConnect() (db *sql.DB) {
	dbDriver := "mysql"
	dbUser := "cerber"
	dbPass := "OU-Shit"
	dbPass := "Girlsaretastybithes"
	dbHost := "127.0.0.1"
	dbName := "shwork"
	db, err := sql.Open(dbDriver, dbUser +":"+ dbPass +"@tcp("+ dbHost +":"+ dbPort +")/"+ dbName +"?charset=utf8")
	if err != nil {
		return
	}
	return db
}
