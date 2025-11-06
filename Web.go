package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type User struct {
	Login    string `json:"login"`
	Password string `json:"password"`
	Gender   string `json:"gender"`
	Age      int    `json:"age"`
	Balance  int    `json:"balance"`
}

var (
	storage = make(map[string]User)
)

func registerHandler(w http.ResponseWriter, r *http.Request) {
	var user User

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Ошибка при чтении JSON", http.StatusBadRequest)
		return
	}

	hashPass := hashPassword(user.Password)
	storage[hashPass] = user

	fmt.Println("Зарегистрирован:", user.Login)
	fmt.Println("Хэш пароля:", hashPass)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status": "registered",
		"user":   user.Login,
		"key":    hashPass,
	})
}

func hashPassword(password string) string {
	hash := sha256.Sum256([]byte(password))
	return hex.EncodeToString(hash[:])
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Ошибка при чтении JSON", http.StatusBadRequest)
		return
	}

	hashPass := hashPassword(user.Password)

	login, exists := storage[hashPass]
	if !exists || login != user {
		http.Error(w, "Неверный логин или пароль", http.StatusUnauthorized)
		return
	}

	fmt.Println("Успешный вход:", user.Login)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status": "success",
		"msg":    "Вы успешно вошли!",
		"key":    hashPass,
	})
}

func infoHandler(w http.ResponseWriter, r *http.Request) {
	key := r.Header.Get("Key")
	if key == "" {
		http.Error(w, "Введите корректный ключ (header Key)", http.StatusUnauthorized)
		return
	}

	user, exists := storage[key]
	if !exists {
		http.Error(w, "Неверный или несуществующий ключ", http.StatusForbidden)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"Login":   user.Login,
		"Key":     key,
		"Info":    "Ключ действителен. Пользователь найден.",
		"Gender":  user.Gender,
		"Age":     strconv.Itoa(user.Age),
		"Balance": strconv.Itoa(user.Balance),
	})
}

func addInfoHandler(w http.ResponseWriter, r *http.Request) {
	key := r.Header.Get("Key")
	if key == "" {
		http.Error(w, "Нет ключа", http.StatusUnauthorized)
		return
	}

	storedUser, exists := storage[key]
	if !exists {
		http.Error(w, "Неверный ключ", http.StatusForbidden)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&storedUser); err != nil {
		http.Error(w, "Ошибка JSON", http.StatusBadRequest)
		return
	}

	storage[key] = storedUser

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"Info":    "Данные добавлены",
		"Gender":  storedUser.Gender,
		"Age":     strconv.Itoa(storedUser.Age),
		"Balance": strconv.Itoa(storedUser.Balance),
	})

}

func main() {
	http.HandleFunc("/addInfo", addInfoHandler)
	http.HandleFunc("/register", registerHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/info", infoHandler)

	if err := http.ListenAndServe(":8000", nil); err != nil {
		fmt.Println("Ошибка запуска сервера:", err)
	}
}
