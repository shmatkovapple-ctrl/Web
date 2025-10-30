package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
)

type User struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

var (
	storage      = make(map[string]string)
	loginStorage = make(map[string]string)
)

func registerHandler(w http.ResponseWriter, r *http.Request) {
	var user User

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		fmt.Println("Ошибка при обработке json", err)
	}
	data := map[string]string{
		"login":    user.Login,
		"password": user.Password,
	}

	hashReg := hashMap(data)

	storage[user.Login] = hashReg

	fmt.Println("Зарегистрирован:", data)
	fmt.Println("Хэш:", hashReg)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status": "registered",
		"user":   user.Login,
	})

}

func hashMap(data map[string]string) string {
	keys := make([]string, 0, len(data))
	for k := range data {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	combined := ""
	for _, k := range keys {
		combined += fmt.Sprintf("%s=%s;", k, data[k])
	}

	hash := sha256.Sum256([]byte(combined))
	return hex.EncodeToString(hash[:])
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	var user User

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		fmt.Println("Ошибка при обработке json", err)
	}
	storedHash, exists := storage[user.Login]
	if !exists {
		http.Error(w, "Пользователь не найден", http.StatusNotFound)
		return
	}

	data := map[string]string{
		"login":    user.Login,
		"password": user.Password,
	}
	hashLog := hashMap(data)

	fmt.Println("Попытка входа:", user.Login)
	fmt.Println("Хэш логина:", hashLog)

	if storedHash == hashLog {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status": "success",
			"msg":    "Вы успешно вошли!",
		})
	} else {
		http.Error(w, "Неверный логин или пароль", http.StatusUnauthorized)
	}
}

func main() {
	http.HandleFunc("/register", registerHandler)
	http.HandleFunc("/login", loginHandler)

	if err := http.ListenAndServe(":8000", nil); err != nil {
		fmt.Println("Произошла ошибка при запуске сервера", err)
	}
}
