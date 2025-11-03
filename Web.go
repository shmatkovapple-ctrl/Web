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
	storage = make(map[string]string)
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

	hash := sha256.Sum256([]byte(keys[0] + data["login"] + data["password"]))
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
			"key":    hashLog,
		})
	} else {
		http.Error(w, "Неверный логин или пароль", http.StatusUnauthorized)
	}
}

func infoHandler(w http.ResponseWriter, r *http.Request) {
	key := r.Header.Get("Key")
	if key == "" {
		http.Error(w, "Введите корректный ключ, поле не может быть пустым", http.StatusUnauthorized)
		return
	}

	var foundUser string
	for login, storedHash := range storage {
		if storedHash == key {
			foundUser = login
			break
		}
	}

	if foundUser == "" {
		http.Error(w, "Неправильный или несуществующий ключ", http.StatusForbidden)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"Login": foundUser,
		"Key":   key,
		"Info":  "Ключ действителен. Пользователь найден.",
	})
}

func main() {
	http.HandleFunc("/info", infoHandler)
	http.HandleFunc("/register", registerHandler)
	http.HandleFunc("/login", loginHandler)

	if err := http.ListenAndServe(":8000", nil); err != nil {
		fmt.Println("Произошла ошибка при запуске сервера", err)
	}
}
