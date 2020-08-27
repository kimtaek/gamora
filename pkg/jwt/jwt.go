package jwt

import (
	"fmt"
	"github.com/caarlos0/env"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"strings"
)

type Configure struct {
	JwtSecret string `env:"JWT_SECRET" envDefault:"Y8TzKMOY5QACbF71m..."`
	Ttl       int64  `env:"JWT_TTL" envDefault:"1"`   // 1 hour
	Trt       int64  `env:"JWT_TRT" envDefault:"168"` // 7 days
}

var Config Configure

func Setup() {
	_ = env.Parse(&Config)
}

func ExtractToken(c *gin.Context) *string {
	bearToken := c.Request.Header.Get("Authorization")
	strArr := strings.Split(bearToken, " ")

	if len(strArr) == 2 {
		return &strArr[1]
	}
	if len(strArr) == 1 {
		return &strArr[0]
	}
	return nil
}

func Verify(token string) (*jwt.Token, error) {
	return jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(Config.JwtSecret), nil
	})
}
