package httpd

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/alexedwards/argon2id"
	"github.com/go-chi/chi"
	"neighbourlink-api/db"
	"neighbourlink-api/types"
	"net/http"
	"regexp"
	"time"
)

func CreateUser(w http.ResponseWriter, r *http.Request, db *db.DB) {
	//  Used to	measure how long the function takes to run (part 1 of 2)
	//	t1 := time.Now()

	var u types.User

	err := json.NewDecoder(r.Body).Decode(&u)

	PasswordHash, err := argon2id.CreateHash(u.Password, argon2id.DefaultParams)
	u.Password = PasswordHash

	re := regexp.MustCompile(`^[a-zA-Z0-9_.+-]+@[a-zA-Z0-9-]+\.[a-zA-Z0-9-.]+$`)
	if re.MatchString(u.Email) == false{
		w.WriteHeader(422)
		http.Error(w, fmt.Sprintf("Email Is Not Valid: %s", u.Email), http.StatusBadRequest)
		fmt.Printf("User tried to make an account with a Email that does not match regex. Email:%s | Username:%s", u.Email, u.Email)
		return
	}

	u.Reputation = 0
	u.Permissions = `Basic`

	_, err = db.AddUser(u)
	if err != nil{
		//user exists
		http.Error(w, fmt.Sprintf("Username exists :%s", u.Username), http.StatusConflict)
		fmt.Printf("User tried to make an account with a Email that already exists. Email:%s | Username:%s", u.Email, u.Username)

		return
	}
	w.WriteHeader(http.StatusCreated)
	return
	//  Used to	measure how long the function takes to run (part 2 of 2)
	//	t2 := time.Now()
	//	diff := t2.Sub(t1)
	//	fmt.Println(diff)

}

func RetrieveUser(w http.ResponseWriter, r *http.Request, db *db.DB) {
	//  Used to	measure how long the function takes to run (part 1 of 2)
	//	t1 := time.Now()


	Username := chi.URLParam(r, "Username")
	var u types.User
	switch Username {
	case "":
		u = types.User{UserID: JWTUserID(r)}
		u, _ = db.GetUserByID(u)

		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(u.Data())
		return

	default:
		u = types.User{Username: Username}
		u, _ = db.GetUserByUsername(u)

		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(u.Data())
		return
	}
}


func LoginUser(w http.ResponseWriter, r *http.Request, db *db.DB, secretKey []byte) {

	var u types.User
	err := json.NewDecoder(r.Body).Decode(&u)

	UserSearch,err := db.GetUserByUsername(u)

	if err != nil{
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}

	check, err := argon2id.ComparePasswordAndHash(u.Password, UserSearch.Password)

	if err != nil {
		fmt.Println(err)
	}

	if check == true {
		details ,_:= json.Marshal(UserSearch) // returns in []byte form
		B64Details := base64.StdEncoding.EncodeToString(details)  // convert to Base64 as there are special characters that are not allowed to be in cookies
		expiry := time.Now().Add(time.Hour * 24 * 7)

		// fmt.Println(string(details))  // need to cast to string to print

		detailsCookie := http.Cookie{
			Name:    "Account",
			Value:   B64Details,
			Expires: expiry,
			Path: "/",
		}

		http.SetCookie(w, &detailsCookie)
		SetJWTcookie(w, UserSearch, secretKey)
		return
	} else {
		http.Error(w, "Not Logged In!", 401)
		return
	}


}

