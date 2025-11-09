package auth

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

// струкрура для реквеста при аутенификации
type SignReq struct {
	Password string `json:"password"`
}

// структура для ответа при аутетификации
type SignRes struct {
	Token string `json:"token,omitempty"`
	Error string `json:"error,omitempty"`
}

// формируем переменную секретного ключа для подписи токена //openssl rand -base64 32
var Secret = []byte("wLJoR/JJjaI/+zGpZqqpoxqU1R0d9hhd/GrogEW5qx4=")

// signinHandler обрабатывает запрос на аутентификацию пользователя
func SigninHandler(w http.ResponseWriter, r *http.Request) {
	//Декодировка  JSON-тела запроса в структуру SignRequest
	var req SignReq
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	//Сравнение введённого пользователем пароля с паролем в переменной окружения TODO_PASSWORD
	if req.Password != os.Getenv("TODO_PASSWORD") {
		resp := SignRes{Error: "Неверный пароль"}
		json.NewEncoder(w).Encode(resp)
		return
	}

	//Формируется JWT-токен с хэшем пароля в полезной нагрузке
	passwordHash := sha256.Sum256([]byte(req.Password))
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"passwordHash": fmt.Sprintf("%x", passwordHash),
	})

	// Подпись токена секретным ключом
	tokenString, err := token.SignedString(Secret)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	//Возвращаем токен в поле token JSON-объекта
	resp := SignRes{Token: tokenString}
	json.NewEncoder(w).Encode(resp)
}
