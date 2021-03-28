package httpd

import (
	"fmt"
	"github.com/gorilla/websocket"
	"neighbourlink-api/db"
	"neighbourlink-api/httpd/middleware"
	"net/http"
)

// using default options create websocket upgrader
var upgrader = websocket.Upgrader{CheckOrigin:func(r *http.Request) bool { return true }}

// endpoint to receive and store websocket connections
func serveWs(w http.ResponseWriter, r *http.Request, db *db.DB) {
	// get id from the cookie stored in the jwt cookie
	JWTID := middleware.JWTUserID(r)

	// if the jwt cookie is 0, the user is unauthorised and thus an error is returned
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