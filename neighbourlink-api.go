package main

import (
	"neighbourlink-api/db"
	"neighbourlink-api/httpd"
)


func main(){
	db := db.InitDB()
	defer db.Close()

	s := httpd.New(db)
	s.Start()
}