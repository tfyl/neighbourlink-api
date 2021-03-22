package admin

import (
	"encoding/json"
	"github.com/go-chi/chi"
	"neighbourlink-api/db"
	"neighbourlink-api/httpd/middleware"
	"neighbourlink-api/types"
	"net/http"
	"strconv"
)

func RetrieveAllUsers(w http.ResponseWriter, r *http.Request, db *db.DB) {

	permissionLevel := middleware.JWTPermission(r) // get permission level from middleware

	// basic auth users should not be able to run this, only "Admin" and "Mod"
	if permissionLevel != `Admin` {
		http.Error(w, "Non Admin Auth Users are not allowed", http.StatusUnauthorized)
		return
	}

	users, _ := db.GetUserAll()  // get an array of all users from DB (db, is a class)

	uStack := types.NewStack()   // initialise class for stack

	for _,u := range users{
		uStack.Push(u.Data())    // save the sanitised version of user (without password hash)
	}

	stackLen := uStack.Len()     // method to return length (int)

	uPriorityQueue := types.NewPQueue()           // initialise class for priority queue
	for i:=0;i<stackLen;i++{
		user := uStack.Pop().(types.User)         // pop user from the top of the stack and save in temp variable user
		uPriorityQueue.Push(user,user.Reputation) // add user to priority and use reputation for value of priority
		//fmt.Println("Pushing:", user)
	}

	var orderedUsers []types.User                 // create variable "orderedUsers" that is a array of custom type "user"

	queueLen := uPriorityQueue.Len()              // store length for priority queue
	for i:=0;i<queueLen;i++{                      // iterate through priority queue
		user := uPriorityQueue.Pop().(types.User) // pop user from front of the queue and save in temp variable user
		//fmt.Println("Popping:", user)
		orderedUsers = append(orderedUsers,user)  // append user to array of sorted users
	}
	//fmt.Println(orderedUsers)
	w.WriteHeader(http.StatusOK)                  // set status code 200 (Status OK)
	_ = json.NewEncoder(w).Encode(orderedUsers)   // write array orderedUsers in JSON form to return
	return


}


func DeleteUser(w http.ResponseWriter, r *http.Request, db *db.DB) {
	permissionLevel := middleware.JWTPermission(r) // get permission level from middleware

	// basic auth users should not be able to run this, only "Admin" and "Mod"
	if permissionLevel != `Admin` {
		http.Error(w, "Non Admin Auth Users are not allowed", http.StatusUnauthorized)
		return
	}

	UserIDstr := chi.URLParam(r, "UserID") // gets UserID taken from the url placeholder
	UserID , err := strconv.Atoi(UserIDstr) // converts UserID from the into a integer
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = db.DeleteUser(types.User{UserID:UserID})  // delete user from DB (db, is a class)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)                  // set status code 200 (Status OK)
	return


}