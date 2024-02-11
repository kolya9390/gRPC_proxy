package jwt_token

import (
	"fmt"
	"log"

	"github.com/golang-jwt/jwt/v5"
	"github.com/kolya9390/gRPC_GeoProvider/client_Proxy/config"
)

func ValidateToken(tokenString string) (*jwt.Token, error) {
    cfg := config.NewAppConf("client_app/.env")

    // Парсим токен с использованием функции Parse из библиотеки jwt
    token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        // Проверяем, что алгоритм подписи соответствует ожидаемому
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			
        }
        // Возвращаем секретный ключ для проверки подписи
        return []byte(cfg.Token.AccessSecret), nil
    })

	log.Println(cfg.Token)
	log.Println(token)
    if err != nil {
		log.Println("ERR1",err)
        return nil, err
    }

    // Проверяем, что токен валиден
    if !token.Valid {
		log.Println("ERR2")
        return nil, fmt.Errorf("token is not valid")
    }

    return token, nil
}