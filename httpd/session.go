package httpd

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/alexedwards/argon2id"
	"neighbourlink-api/db"
	"neighbourlink-api/httpd/middleware"
	"neighbourlink-api/types"
	"net/http"
	"time"
)

// create session / login user
func CreateSession(w http.ResponseWriter, r *http.Request, db *db.DB, secretKey []byte) {

	var u types.User
	// decode the data from user into user object
	err := json.NewDecoder(r.Body).Decode(&u)

	// get user details with the user id
	UserSearch,err := db.GetUserByUsername(u)

	if err != nil{
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	// compare hash of password and the hash stored
	check, err := argon2id.ComparePasswordAndHash(u.Password, UserSearch.Password)

	if err != nil {
		fmt.Println(err)
	}

	if check == true {
		// returns in []byte form
		details ,_:= json.Marshal(UserSearch.Data())
		// convert to Base64 as there are special characters that are not allowed to be in cookies
		B64Details := base64.StdEncoding.EncodeToString(details)
		// sets expiry 7 days from now
		expiry := time.Now().Add(time.Hour * 24 * 7)
		// creates cookie to store account in json base64 encoded format
		detailsCookie := http.Cookie{
			Name:    "Account",
			Value:   B64Details,
			Expires: expiry,
			Path: "/",
		}
		// set cookie client side
		http.SetCookie(w, &detailsCookie)
		// set jwt cookie for authentication
		middleware.SetJWTcookie(w, UserSearch, secretKey)
		return
	} else {
		http.Error(w, "Not Logged In!", 401)
		return
	}


}

