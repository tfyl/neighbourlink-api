package db

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

const connStr string = `host=ec2-54-247-125-38.eu-west-1.compute.amazonaws.com
				user=ofeewfrhqunjte
				password=027b9aec6e5217b1d8a64183bfa85e4af64b22b43c366011339ad42856b1de50
				port=5432
				dbname=d27bhqu6jrkt4h
				sslmode=require`

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
	user_id     	      SERIAL REFERENCES user_detail(user_id),
	post_time 			  timestamp,
	post_title            varchar ,
	post_description      varchar,
	post_urgency          varchar
);

CREATE TABLE IF NOT EXISTS post_comment(
	comment_id       SERIAL PRIMARY KEY,
	post_id          SERIAL REFERENCES post(post_id),
	user_id          SERIAL REFERENCES user_detail(user_id) ,
	comment_message  varchar
);`


func InitDB() *DB {
	var db, err = sqlx.Open("postgres", connStr)
	if err != nil{
		fmt.Println(err)
	}
	_ , err= db.Exec(schema)
	if err != nil{
		fmt.Println(err)
	}

	return &DB{*db}

}


type DB struct {
	sqlx.DB
}

// remove all tx for singular select