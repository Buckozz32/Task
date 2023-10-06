package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)


func refresh(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("rt")
	rtString := c.Value
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	rtClaims, err := getClaims(c)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	c, err = r.Cookie("at")
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	atClaims, err := getClaims(c)
	if err != nil && err.Error() != "Token is not valid" {
		fmt.Println("xd not vaild")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

d
	if rtClaims.UserID != atClaims.UserID || rtClaims.Id != atClaims.Id {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}


	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(
		connString,
	))
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}()

	var session mongo.Session
	if session, err = client.StartSession(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err = session.StartTransaction(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err = mongo.WithSession(ctx, session, func(sc mongo.SessionContext) error {
		collection := client.Database("goauth").Collection("users")
		findFilter := bson.M{"guid": rtClaims.UserID}
		var result User

		err = collection.FindOne(ctx, findFilter).Decode(&result)
		if err != nil {
			fmt.Printf("collection find err: %v\n", err)
			return err
		}

		var rtExists bool
		for _, hash := range result.Rts {
			
			err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(reverse(rtString)))
			if err == nil {
				rtExists = true
				break
			}
		}
		if !rtExists {
			w.WriteHeader(http.StatusUnauthorized)
			return errors.New("Not found")
		}
		return nil
	}); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	session.EndSession(ctx)


	atExpiration := time.Now().Add(1 * time.Minute)
	atClaims = &сlaims{
		UserID: atClaims.UserID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: atExpiration.Unix(),
		
			Id: rtClaims.Id,
		}
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS512, atClaims)
	atSigned, err := accessToken.SignedString(jwtKey)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "at",
		Value:    atSigned,
		Expires:  atExpiration,
		HttpOnly: true,
	})
}
