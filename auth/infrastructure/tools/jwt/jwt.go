package jwt

import (
	"auth/config"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type User struct {
	ID       int64
	Email    string
}

func NewToken(user User) (string, error) {

	cfg := config.NewAppConf("server_app/.env")

	token := jwt.New(jwt.SigningMethodHS512)

	claims := token.Claims.(jwt.MapClaims)

	claims["uid"] = user.ID
	claims["email"] = user.Email
	claims["exp"] = time.Now().Add(cfg.Token.AccessTTL).Unix()

	tokenString, err := token.SignedString([]byte(cfg.Token.AccessSecret))

	if err != nil {
		return "", fmt.Errorf("new token error : %s",err)
	}

	return tokenString,err

}
