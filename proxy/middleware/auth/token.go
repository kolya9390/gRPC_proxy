package auth_token_midw

import (
	"net/http"

	jwt_token "github.com/kolya9390/gRPC_GeoProvider/client_Proxy/infrastructure/jwt"
)

func TokenAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		// Проверяем наличие токена в заголовке запроса
		if token == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Проверяем валидность токена
		_, err := jwt_token.ValidateToken(token)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Прошла аутентификация, передаем управление следующему обработчику
		next.ServeHTTP(w, r)
	})
}
