package utils

import (
	"bigmind/xcheck-be/config"
	"bigmind/xcheck-be/internal/models"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func GenerateToken(user *models.User) (map[string]string, error) {
	tokenLifespan := config.GetConfig().GetInt("AUTH_JWT_EXPIRE")

	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["username"] = user.Username
	claims["sub"] = user.ID
	claims["exp"] = time.Now().Add(time.Hour * time.Duration(tokenLifespan)).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(config.GetConfig().GetString("AUTH_JWT_SECRET")))
	if err != nil {
		return nil, err
	}

	rtClaims := jwt.MapClaims{}
	rtClaims["sub"] = user.ID
	rtClaims["exp"] = time.Now().Add(time.Hour * 24).Unix()
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	rt, err := refreshToken.SignedString([]byte(config.GetConfig().GetString("AUTH_JWT_SECRET")))
	if err != nil {
		return nil, err
	}

	return map[string]string{
		"token":         t,
		"refresh_token": rt,
	}, nil

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

func ExtractTokenID(c *gin.Context) (int64, string, error) {
	tokenString := ExtractToken(c)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(config.GetConfig().GetString("AUTH_JWT_SECRET")), nil
	})
	if err != nil {
		return 0, "", err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		log.Println(claims)
		uname := claims["username"].(string)
		uid := claims["sub"].(float64)
		return int64(uid), uname, nil
	}
	return 0, "", nil
}
