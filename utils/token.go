package utils

import (
	"bigmind/xcheck-be/config"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func GenerateToken(username string) (string, error) {
	tokenLifespan := config.GetConfig().GetInt("AUTH_JWT_EXPIRE")

	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["username"] = username
	claims["exp"] = time.Now().Add(time.Hour * time.Duration(tokenLifespan)).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(config.GetConfig().GetString("AUTH_JWT_SECRET")))

}

func TokenValid(c *gin.Context) error {
	tokenString := ExtractToken(c)

	// log.Println("tokenString: ", tokenString)
	// log.Println("jwtSecret: ", config.GetConfig().GetString("auth.jwtSecret"))

	if tokenString == "" {
		return errors.New("Unauthorized")
	}

	_, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(config.GetConfig().GetString("AUTH_JWT_SECRET")), nil
	})
	if err != nil {
		return err
	}
	return nil
}

func ExtractToken(c *gin.Context) string {
	token := c.Query("token")
	if token != "" {
		return token
	}
	bearerToken := c.Request.Header.Get("Authorization")
	if len(strings.Split(bearerToken, " ")) == 2 {
		return strings.Split(bearerToken, " ")[1]
	}
	return ""
}

func ExtractTokenID(c *gin.Context) (string, error) {

	tokenString := ExtractToken(c)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(config.GetConfig().GetString("AUTH_JWT_SECRET")), nil
	})
	if err != nil {
		return "", err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		uname := claims["username"].(string)
		return uname, nil
	}
	return "", nil
}
