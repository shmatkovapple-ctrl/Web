// 1. Подключение БД
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

// 1. Подключение БД
// 2. Сделать структуру, хэндлер регистрации

var (
	db     *pgxpool.Pool
	jwtKey []byte
)

type Users struct {
	Id        int       `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	BirthDate time.Time `json:"birth_date"`
	CreatedAt time.Time `json:"created_at"`
	Money     int       `json:"money"`
}

type LoginReq struct {
	FirstName string `json:"first_name"`
}

func generateJWT(userID int, firstName string) (string, error) {
	claims := jwt.MapClaims{
		"user_id":    userID,
		"first_name": firstName,
		"exp":        time.Now().Add(24 * time.Hour).Unix(), // токен действителен 24 часа
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

func verifyJWT(r *http.Request) (jwt.MapClaims, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return nil, fmt.Errorf("нет токена в заголовке Authorization")
	}

	var tokenStr string
	fmt.Sscanf(authHeader, "Bearer %s", &tokenStr)

	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("неверный метод подписи")
		}
		return jwtKey, nil
	})

	if err != nil || !token.Valid {
		return nil, fmt.Errorf("недействительный токен: %v", err)
	}

	claims := token.Claims.(jwt.MapClaims)
	return claims, nil
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	type Reg struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		BirthDate string `json:"birth_date"`
	}

	var reg Reg
	if err := json.NewDecoder(r.Body).Decode(&reg); err != nil {
		http.Error(w, "Ошибка записи в регистрации", 400)
	}

	var id int

	query := `INSERT INTO users (first_name, last_name, birth_date) 
	VALUES ($1, $2, $3)
	RETURNING id`

	err := db.QueryRow(context.Background(), query,
		reg.FirstName, reg.LastName, reg.BirthDate).Scan(&id)

	if err != nil {
		http.Error(w, "Ошибка при добавлении в БД: "+err.Error(), 500)
		return
	}

	_, err = db.Exec(
		context.Background(),
		`INSERT INTO wallets (user_id) VALUES ($1)`,
		id,
	)
	if err != nil {
		http.Error(w, "Ошибка создания кошелька", 500)
		return
	}

	token, err := generateJWT(id, reg.FirstName)
	if err != nil {
		http.Error(w, "Ошибка генерации токена", 500)
		return
	}

	fmt.Fprintf(w, `{"status":"ok","id":%d}`, id, token)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	type LoginReq struct {
		FirstName string `json:"first_name"`
	}

	var req LoginReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Ошибка JSON", http.StatusBadRequest)
		return
	}

	claims, err := verifyJWT(r)
	if err != nil {
		http.Error(w, "Недействительный токен: "+err.Error(), http.StatusUnauthorized)
		return
	}

	tokenFirstName, ok := claims["first_name"].(string)
	if !ok || tokenFirstName == "" {
		http.Error(w, "Некорректный токен: нет имени", http.StatusUnauthorized)
		return
	}

	authSuccess := false
	if tokenFirstName == req.FirstName {
		authSuccess = true
	} else {
		authSuccess = true
	}

	if !authSuccess {
		http.Error(w, "Отказ в авторизации", http.StatusForbidden)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status": "success",
		"msg":    "Авторизация успешна",
		"name":   req.FirstName,
		"token":  r.Header.Get("Authorization"),
	})
}

func usersHandler(w http.ResponseWriter, r *http.Request) {

	claims, err := verifyJWT(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	fmt.Println("JWT claims:", claims)

	var users []Users

	rows, err := db.Query(context.Background(), `SELECT id, first_name, last_name, birth_date, created_at, money FROM users`)
	if err != nil {
		http.Error(w, "Ошибка при запросе к БД: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var u Users
		err := rows.Scan(&u.Id, &u.FirstName, &u.LastName, &u.BirthDate, &u.CreatedAt, &u.Money)
		if err != nil {
			http.Error(w, "Ошибка при чтении данных: "+err.Error(), http.StatusInternalServerError)
			return
		}
		users = append(users, u)
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(users); err != nil {
		http.Error(w, "Ошибка при кодировании JSON: "+err.Error(), http.StatusInternalServerError)
	}
}

func insertMoneyHandler(w http.ResponseWriter, r *http.Request) {
	claims, err := verifyJWT(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}
	userID := int(userIDFloat)

	type Req struct {
		Amount int `json:"amount"`
	}

	var req Req
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Bad JSON", http.StatusBadRequest)
		return
	}

	if req.Amount <= 0 {
		http.Error(w, "Invalid amount", http.StatusBadRequest)
		return
	}

	var balance int
	err = db.QueryRow(
		context.Background(),
		`UPDATE wallets
		 SET balance = balance + $1
		 WHERE user_id = $2
		 RETURNING balance`,
		req.Amount, userID,
	).Scan(&balance)

	if err != nil {
		http.Error(w, "DB error", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]any{
		"status":  "ok",
		"balance": balance,
	})
}

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("⚠️  .env файл не найден или ошибка загрузки")
	}

	jwtKey = []byte(os.Getenv("JWT_SECRET"))

	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	connStr := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		dbUser, dbPass, dbHost, dbPort, dbName,
	)
	if err != nil {
		fmt.Println("Ошибка подключения:", err)
		return
	}
	db, err = pgxpool.New(context.Background(), connStr)
	if err != nil {
		fmt.Println("Ошибка подключения:", err)
		return
	}
	defer db.Close()

	fmt.Println("Успешное подключение к PostgreSQL!")

	http.HandleFunc("/register", registerHandler)
	http.HandleFunc("/users", usersHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/wallet/deposit", insertMoneyHandler)

	fmt.Println("Server started at :8000")
	if err := http.ListenAndServe(":8000", nil); err != nil {
		fmt.Println("Ошибка запуска сервера:", err)
	}
}
