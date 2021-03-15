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


	claims, _ := sjwt.ToClaims(u.Data())
	expiry := time.Now().Add(time.Hour * 24 * 7)
	claims.SetExpiresAt(expiry)

	JWTstring := claims.Generate(secretKey)

	JWTcookie := http.Cookie{
		Name:    "JWT",
		Value:   JWTstring,
		Expires: expiry,
		Path: "/",
	}

	http.SetCookie(w, &JWTcookie)
	_, _ = w.Write([]byte("Login Succesful"))
}


func JWTAuthMiddleware(next http.Handler,secretKey []byte) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		JWTcookie, err := r.Cookie("JWT")
		if err != nil{
			fmt.Println(err)
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		claims, err := sjwt.Parse(JWTcookie.Value)
		if err != nil {
			fmt.Println(err)
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		var user types.User

		err = claims.ToStruct(&user)
		if err != nil {
			fmt.Println(err)
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		check := sjwt.Verify(JWTcookie.Value, secretKey)
		if !check {
			fmt.Println(err)
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		r = r.WithContext(context.WithValue(r.Context(),"Permission",user.Permissions))
		r = r.WithContext(context.WithValue(r.Context(),"UserID",user.UserID))

		next.ServeHTTP(w, r)
	})
}

func JWTUserID(r *http.Request) int {
	if r.Context().Value("UserID") == nil{
		fmt.Println("No ctx value set for UserID")
		return 0
	}
	JWTID := r.Context().Value("UserID").(int)
	return JWTID
}

func JWTPermission(r *http.Request) string {
	if r.Context().Value("Permission") == nil{
		fmt.Println("No ctx value set for Permission")
		return ""
	}
	JWTperm := r.Context().Value("Permission").(string)
	return JWTperm
}