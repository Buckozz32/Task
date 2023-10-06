package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/dgrijalva/jwt-go"
)

type сlaims struct {
	UserID string `json:"uid"`
	jwt.StandardClaims
}

func getClaims(c *http.Cookie) (*сlaims, error) {
	tknStr := c.Value
	claims := &сlaims{}


	tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
		
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			
			return nil, fmt.Errorf("Неожиданный алгоритм подписи: %v", token.Header["alg"])
		}

		return jwtKey, nil
	})
	if err != nil {
		//fmt.Printf("parse err: %v", err)
		return nil, err
	}
	if !tkn.Valid {
		//fmt.Println("tk not valid")
		return claims, errors.New("Token is not valid")
	}

	return claims, nil
}

func reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}
