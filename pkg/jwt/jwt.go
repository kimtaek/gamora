package jwt

import (
	"fmt"
	"github.com/caarlos0/env"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"strings"
)

// Configure config for jwt
type Configure struct {
	JwtSecret string `env:"JWT_SECRET" envDefault:"Y8TzKMOY5QACbF71m..."`
	TTL       int64  `env:"JWT_TTL" envDefault:"1"`   // 1 hour
	TRT       int64  `env:"JWT_TRT" envDefault:"168"` // 7 days
}

// Config global defined jwt config
var Config Configure

// Setup init jwt
func Setup() {
	_ = env.Parse(&Config)
}

// ExtractToken return jwt token string or nil
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

// Verify verify token sign
func Verify(token string) (*jwt.Token, error) {
	return jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(Config.JwtSecret), nil
	})
}
