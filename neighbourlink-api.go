package main

import (
	"neighbourlink-api/db"
	"neighbourlink-api/httpd"
)


func main(){
	// initialises database
	db := db.InitDB()
	// defer closing the database (close when function is ended)
	defer db.Close()

	// creates httpd object with database
	s := httpd.New(db)
	// starts http server
	s.Start()
}