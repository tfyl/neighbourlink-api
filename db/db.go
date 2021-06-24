package db

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"neighbourlink-api/types"
)

const connStr string = `host=
				user=ofeewfrhqunjte
				password=
				port=
				dbname=
				sslmode=`

//noinspection SqlNoDataSourceInspection
var schema = `
CREATE TABLE IF NOT EXISTS user_detail(
	user_id SERIAL PRIMARY KEY,
	username       varchar UNIQUE,
	email 	       varchar
);

CREATE TABLE IF NOT EXISTS user_auth(
	user_id SERIAL PRIMARY KEY,
	password       varchar,
	permissions    varchar,

	FOREIGN KEY(user_id) 
	  REFERENCES user_detail(user_id)
	  	ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS user_attribute(
	user_id SERIAL PRIMARY KEY,
	local_area       varchar,
	reputation       INT,
	FOREIGN KEY(user_id) 
	  REFERENCES user_detail(user_id)
	  	ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS post(
	post_id               SERIAL PRIMARY KEY,
	user_id     	      SERIAL,
	post_time 			  timestamp,
	post_title            varchar ,
	post_description      varchar,
	post_urgency          varchar,
	FOREIGN KEY(user_id) 
	  REFERENCES user_detail(user_id)
	  	ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS post_comment(
	comment_id       SERIAL PRIMARY KEY,
	post_id          SERIAL,
	user_id          SERIAL,
	comment_message  varchar,
	FOREIGN KEY(post_id) 
	  REFERENCES post(post_id)
	  	ON DELETE CASCADE,
	FOREIGN KEY(user_id) 
	  REFERENCES user_detail(user_id)
	  	ON DELETE CASCADE
);`
// Fully Normalised Schema for PostgreSQL table


func InitDB() *DB {
	// initialises database
	var db, err = sqlx.Open("postgres", connStr) // opens connection to the postgres server
	if err != nil{
		fmt.Println(err)
	}
	_ , err= db.Exec(schema) // executes the schema (only creates table "IF NOT EXISTS")
	if err != nil{
		fmt.Println(err)
	}

	webSocketMap := types.NewWebSocketMap() // creates a map/dictionary of all websockets

	return &DB{DB:*db,Websocket:webSocketMap} // returns the DB type
}


type DB struct {
	sqlx.DB // DB structure contains a connection to the database
	Websocket types.WebSocketMapStruct // contains a custom structure to hold the websocket connections
}

