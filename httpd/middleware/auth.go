package middleware

import (
	"context"
	"fmt"
	"github.com/brianvoe/sjwt"
	"neighbourlink-api/types"
	"net/http"
	"time"
)

func SetJWTcookie(w http.ResponseWriter, u types.User , secretKey []byte) {
	// set JWT cookie

	claims, _ := sjwt.ToClaims(u.Data())
	// adds user claims (u.Data returns user class without password hash)
	expiry := time.Now().Add(time.Hour * 24 * 7)
	// sets cookie expiry to 7 days in the future
	claims.SetExpiresAt(expiry) // sets claims to expire in 7 days

	JWTstring := claims.Generate(secretKey)
	// gets the cookie data in the form of a string
	// creates the cookie struct/object
	JWTcookie := http.Cookie{
		Name:    "JWT", // name of the cookie
		Value:   JWTstring, // value stored within cookie (data)
		Expires: expiry, // expiry date
		Path: "/", // the path for which the cookie is valid ("/" means all paths)
	}

	http.SetCookie(w, &JWTcookie)
	_, _ = w.Write([]byte("Login Succesful"))
}


func JWTAuthMiddleware(next http.Handler,secretKey []byte) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		JWTcookie, err := r.Cookie("JWT")
		// gets cookie called JWT that was sent along with the request
		if err != nil{
			fmt.Println(err)
			http.Error(w, err.Error(), http.StatusUnauthorized)
			// if cookie doesn't exist, return error 401 (Unauthorised)
			return
		}

		claims, err := sjwt.Parse(JWTcookie.Value)
		// Parse the claims from the value of the cookie
		if err != nil {
			fmt.Println(err)
			http.Error(w, err.Error(), http.StatusUnauthorized)
			// if the claims are invalid, return error 401 unauthorised
			return
		}

		var user types.User
		// instantiates an object to assign the claims to

		err = claims.ToStruct(&user) // assigns claims to user object
		if err != nil {
			fmt.Println(err)
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		check := sjwt.Verify(JWTcookie.Value, secretKey)
		// verify the JWT cookie with the secret key to check if the data has been tampered with
		if !check {
			fmt.Println(err)
			http.Error(w, err.Error(), http.StatusUnauthorized)
			// if the data has been tampered with, return 401 error
			return
		}

		// pass along context with value (so the following request can simply extract the permissions and userid value)
		r = r.WithContext(context.WithValue(r.Context(),"Permission",user.Permissions))
		r = r.WithContext(context.WithValue(r.Context(),"UserID",user.UserID))

		// execute the code for the request
		next.ServeHTTP(w, r)
	})
}

func JWTUserID(r *http.Request) int {
	// get user id from the JWT token and return it in an int form
	if r.Context().Value("UserID") == nil{
		fmt.Println("No ctx value set for UserID")
		return 0
	}
	JWTID := r.Context().Value("UserID").(int)
	return JWTID
}

func JWTPermission(r *http.Request) string {
	// get user permission level from the JWT token and return it in an string form
	if r.Context().Value("Permission") == nil{
		fmt.Println("No ctx value set for Permission")
		return ""
	}
	JWTperm := r.Context().Value("Permission").(string)
	return JWTperm
}