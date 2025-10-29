package main

import (
	"fmt"
	"net/http"
)

func main() {
	if err := http.ListenAndServe(":8000", http.FileServer(http.Dir("."))); err != nil {
		fmt.Println("Произошла ошибка при запуске сервера", err)
	}
}
