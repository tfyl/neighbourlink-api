package httpd

import (
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
)

func CreateUser(w http.ResponseWriter, r *http.Request, db *db.DB) {
	//  Used to	measure how long the function takes to run (part 1 of 2)
	//	t1 := time.Now()

	var u types.User

	err := json.NewDecoder(r.Body).Decode(&u)

	PasswordHash, err := argon2id.CreateHash(u.Password, argon2id.DefaultParams)
	u.Password = PasswordHash

	// compiles regex to check email
	re := regexp.MustCompile(`^[a-zA-Z0-9_.+-]+@[a-zA-Z0-9-]+\.[a-zA-Z0-9-.]+$`)
	// checks if email matches compiled regex
	if re.MatchString(u.Email) == false{
		// returns error if the input is not valid
		w.WriteHeader(422)
		http.Error(w, fmt.Sprintf("Email Is Not Valid: %s", u.Email), http.StatusBadRequest)
		fmt.Printf("User tried to make an account with a Email that does not match regex. Email:%s | Username:%s", u.Email, u.Email)
		return
	}
	// set user reputation to 1 when signing up
	u.Reputation = 1
	// set user permissions to basic when signing up (i.e. not Admin or Moderator)
	u.Permissions = `Basic`

	// adds user to database
	_, err = db.AddUser(u)
	if err != nil{
		//user exists return error
		http.Error(w, fmt.Sprintf("Username exists :%s", u.Username), http.StatusConflict)
		fmt.Printf("User tried to make an account with a Email that already exists. Email:%s | Username:%s", u.Email, u.Username)

		return
	}

	// sets 301 created header
	w.WriteHeader(http.StatusCreated)
	return
	//  Used to	measure how long the function takes to run (part 2 of 2)
	//	t2 := time.Now()
	//	diff := t2.Sub(t1)
	//	fmt.Println(diff)

}

func RetrieveUser(w http.ResponseWriter, r *http.Request, db *db.DB) {
	// gets url parameter user id
	UserIDs := chi.URLParam(r, "UserID")
	var u types.User
	// if user id is not present it will return the details for the currently logged in user
	switch UserIDs {
	// user id is empty
	case "":
		// get the user id of the currently authenticated user
		u = types.User{UserID: middleware.JWTUserID(r)}
		// get user from database by id
		u, _ = db.GetUserByID(u)
		w.WriteHeader(http.StatusOK)
		// return the user in json form (sanitised)
		_ = json.NewEncoder(w).Encode(u.Data())
		return

	// user id is present
	default:
		// converts user id to a int
		UserID , err := strconv.Atoi(UserIDs)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		u = types.User{UserID: UserID}
		// gets user data
		u, err := db.GetUserByID(u)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
		// return the user in json form (sanitised)
		_ = json.NewEncoder(w).Encode(u.Data())
		return
	}
}


func UpdateUser(w http.ResponseWriter,r *http.Request,db *db.DB){
	JWTID := middleware.JWTUserID(r) // gets user id from signed jwt cookie
	JWTPermission := middleware.JWTPermission(r) // gets permission level of cookie using the middleware that is run for this request

	UserIDstr := chi.URLParam(r, "UserID") // gets UserID taken from the url placeholder
	UserID , err := strconv.Atoi(UserIDstr) // converts UserID from the into a integer
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

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

	originalUser, err := db.GetUserByID(types.User{UserID:UserID}) // defines user object / struct as : "u" ; it gets current record of the user from the user id provided
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	u := originalUser
	fmt.Println(u)

	err = json.NewDecoder(r.Body).Decode(&u) // decoding request body into the user object / struct (overriding the old attributes)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	u.UserID = UserID

	if u.Password == ""{
		u.Password = originalUser.Password
	}
	if JWTPermission != `Admin`{
		u.Reputation = originalUser.Reputation // stop unauthorised `Basic` users changing reputation
		u.LocalArea = originalUser.LocalArea
	}


	u, err = db.UpdateUser(u)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		fmt.Println(err.Error())
		return
	}

	_ = json.NewEncoder(w).Encode(u.Data())

}

