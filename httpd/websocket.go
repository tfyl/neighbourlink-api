package httpd

import (
	"fmt"
	"github.com/gorilla/websocket"
	"neighbourlink-api/db"
	"neighbourlink-api/httpd/middleware"
	"net/http"
)

var upgrader = websocket.Upgrader{CheckOrigin:func(r *http.Request) bool { return true }} // use default options

func serveWs(w http.ResponseWriter, r *http.Request, db *db.DB) {
	JWTID := middleware.JWTUserID(r)

	if JWTID == 0 {
		http.Error(w, "Not authorised", http.StatusUnauthorized)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		fmt.Println(err)
//		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	db.Websocket.AddConn(JWTID,conn)

}