package auth

import (
	"crypto/sha256"
	"diplomaGoSologub/models"
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
)

/*механизм middleware, который принимает http.HandlerFunc и представляет следующий обработчик в цепочке, для
проверки аутентификации для следующих API-запросов:
/api/task — все поддерживаемые HTTP методы;
/api/tasks — получение списка задач;
/api/task/done — запрос на выполнение задачи.
*/

func Auth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		pass := models.PwrdGetEnv()
		if len(pass) > 0 {
			//Если пароль определен, функция получает хэш пароля и преобразует его в строковый формат.
			passwordHash := sha256.Sum256([]byte(pass))
			passwordHashString := fmt.Sprintf("%x", passwordHash)

			// Получаем куку
			var jwtToken string
			cookie, err := r.Cookie("token")
			if err == nil {
				jwtToken = cookie.Value
			}

			//Парсит JWT-токен и проверяет его валидность
			token, err := jwt.Parse(jwtToken, func(token *jwt.Token) (interface{}, error) {
				// Проверяем, что метод подписи токена - HMAC
				//это дополнительная проверка от злоумышленников
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}
				// Возвращаем секретный ключ для проверки подписи токена
				return Secret, nil
			})

			if err != nil {
				http.Error(w, "Authentification required", http.StatusUnauthorized)
				return
			}

			// При валидации JWT токена не забудьте сравнить хэш (или контрольную сумму) текущего пароля и его хэш из токена. Если изменился пароль в переменной окружения, проверка токена для старого пароля не может быть валидной;
			if !token.Valid {
				http.Error(w, "Authentification required", http.StatusUnauthorized)
				return
			}

			//fmt.Println(token)
			// Извлечение  полезные данные из токена и проверяем хэш пароля https://www.jwt.io
			// приводим поле Claims к типу jwt.MapClaims
			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok { // если Сlaims вдруг оказжется другого типа, мы получим панику
				http.Error(w, "Authentification required", http.StatusUnauthorized)
				return
			}

			//fmt.Println(claims)
			//Так как jwt.Claims — словарь вида map[string]inteface{}, используем синтакис получения значения по ключу
			hRaw := claims["passwordHash"]
			h, ok := hRaw.(string)
			if !ok {
				http.Error(w, "Authentification required", http.StatusUnauthorized)
				return
			}

			if h != passwordHashString {
				//В случае ошибки аутентификации следует возвращать ошибку http.StatusUnauthorized (код 401)
				http.Error(w, "Authentification required", http.StatusUnauthorized)
				return
			}
		}

		next(w, r)

	})
}
