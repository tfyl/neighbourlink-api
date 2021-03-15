package httpd

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/alexedwards/argon2id"
	"github.com/go-chi/chi"
	"neighbourlink-api/db"
	"neighbourlink-api/httpd/middleware"
	"neighbourlink-api/types"
	"net/http"
	"regexp"
	"strconv"
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


	UserIDs := chi.URLParam(r, "UserID")
	var u types.User
	switch UserIDs {
	case "":
		u = types.User{UserID: middleware.JWTUserID(r)}
		u, _ = db.GetUserByID(u)

		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(u.Data())
		return

	default:
		UserID , err := strconv.Atoi(UserIDs)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		u = types.User{UserID: UserID}
		u, err := db.GetUserByID(u)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(u.Data())
		return
	}
}


func UpdateUser(w http.ResponseWriter,r *http.Request,db *db.DB){
	JWTID := middleware.JWTUserID(r) // gets user id from signed jwt cookie
	JWTPermission := middleware.JWTPermission(r) // gets permission level of cookie using the middleware that is run for this request

	UserIDstr := chi.URLParam(r, "UserID") // gets UserID taken from the url placeholder
	UserID , err := strconv.Atoi(UserIDstr) // converts UserID from the into a integer

	// Truth table
	// user is if the authenticated user = user being modified
	// admin is if the authenticated user has the permissions of an admin
	// user | admin | desired | and | or  | nand  | nor
	//   1  |   1   |    0    |  1  |  1  |   0   |  0
	//   1  |   0   |    0    |  0  |  1  |   1   |  0
	//   0  |   1   |    0    |  0  |  1  |   1   |  0
	//   0  |   0   |    1    |  0  |  0  |   1   |  1
	//   is desired result :  |  x  |  x  |   x   |  ✔

	if !(UserID == JWTID || JWTPermission == `Admin`) {
		http.Error(w, "user is not authorised to update", http.StatusUnauthorized)
		return
	}

	u, err := db.GetUserByID(types.User{UserID:UserID}) // defines user object / struct as : "u" it gets current record of the user

	err = json.NewDecoder(r.Body).Decode(&u) // decoding request body into the user object / struct
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	u.UserID = UserID

	u, err = db.UpdateUser(u)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_ = json.NewEncoder(w).Encode(u)

}


func LoginUser(w http.ResponseWriter, r *http.Request, db *db.DB, secretKey []byte) {

	var u types.User
	err := json.NewDecoder(r.Body).Decode(&u)

	UserSearch,err := db.GetUserByUsername(u)

	if err != nil{
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	check, err := argon2id.ComparePasswordAndHash(u.Password, UserSearch.Password)

	if err != nil {
		fmt.Println(err)
	}

	if check == true {
		details ,_:= json.Marshal(UserSearch.Data()) // returns in []byte form
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
		middleware.SetJWTcookie(w, UserSearch, secretKey)
		return
	} else {
		http.Error(w, "Not Logged In!", 401)
		return
	}


}

