package main

import (
	"log"
	"net/http"
	"os"
)

/*
	Access, Refresh,  - основные 2 маршрута
	access_module - маршрут для правильности обработки токенов

	env_vars - строки для доступа к БД и ключ для шифрования jwt
	token_operaions - операции с jwt токенами
	user - модель пользователя в БД
*/

func main() {

	http.HandleFunc("/receive", receive)

	http.HandleFunc("/refresh", refresh)
  http.HandleFunc("/access", accessProtectedResource)

	port := os.Getenv("PORT")
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
