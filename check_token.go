ackage main

import (

	"fmt"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	
)

func TokenCheck(c *http.Cookie, err error) (int, *jwt.Claims) {
if err != nil {
	if err == http.ErrNoCookie {
		return http.StatusUnauthorized, nil
		
	}
}
fmt.Println("cookie err")

	
TokenStr := c.Value

Claims := &claims{}

token, err := jwt.ParseWithClaims(TokenStr, Claims, func(token *jwt.Token) (interface{}, error) {
	if _, ok := token.Method(*&jwt.SigningMethodHMAC); !ok {
		fmt.Println("sing method unexpected %v\n", token.Header["algoritm"])
		return nil, fmt.Errorf("Неправильный алгоритм токена: %v", token.Header["algoritm"])
	}
	return jwtKey, nil

})

if err != nil {
	if err == jwt.ErrSignatureInvalid {
		fmt.Println("token sing not valid")
	return http.StatusUnauthorized, nil
	}
}
fmt.Println("parse err: %v", err)
return http.StatusBadRequest, nil


}


func accessProtectedResourse(w http.ResponseWriter, r *http.Request) {
	code, claims := TokenCheck("at")
	if code != http.StatusOK {
		w.WriteHeader(code)
	return
	}
	w.Write([]byte(fmt.Sprintf("Hello, %s!", claims.UserID)))

}
