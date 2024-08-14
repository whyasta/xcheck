package utils

import (
	"bigmind/xcheck-be/config"
	"bigmind/xcheck-be/internal/models"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type AuthDetails struct {
	AuthUuid string
	UserId   uint64
}

type Claims struct {
	UserId uint64 `json:"user_id"`
	jwt.RegisteredClaims
}

func GenerateTokenUuid(authD *AuthDetails) (map[string]string, error) {
	tokenLifespan := config.GetConfig().GetInt("AUTH_JWT_EXPIRE")

	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["auth_uuid"] = authD.AuthUuid
	claims["user_id"] = authD.UserId
	claims["sub"] = authD.UserId
	claims["exp"] = time.Now().Add(time.Hour * time.Duration(tokenLifespan)).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(config.GetConfig().GetString("AUTH_JWT_SECRET")))
	if err != nil {
		return nil, err
	}

	rtClaims := jwt.MapClaims{}
	rtClaims["sub"] = authD.UserId
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

func GenerateToken(user *models.User) (map[string]string, error) {
	tokenLifespan := config.GetConfig().GetInt("AUTH_JWT_EXPIRE")

	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["username"] = user.Username
	claims["auth_uuid"] = uuid.New().String()
	claims["user_id"] = user.ID
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

func TokenValid(c *gin.Context) (error, bool) {
	tokenString := ExtractToken(c)

	// log.Println("tokenString: ", tokenString)
	// log.Println("jwtSecret: ", config.GetConfig().GetString("auth.jwtSecret"))

	if tokenString == "" {
		return errors.New("unauthorized"), false
	}

	// check blacklist or not
	// redisToken, err := redis.String(config.NewRedis().Get().Do("GET", "blacklist_"+tokenString))
	// if redisToken != "" && err == nil {
	// 	return errors.New("unauthorized - token blacklist")
	// }

	_, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(config.GetConfig().GetString("AUTH_JWT_SECRET")), nil
	})
	if err != nil {
		if strings.Contains(err.Error(), "expired") {
			return err, true
		}
		return err, false
	}

	// claims, ok := token.Claims.(jwt.MapClaims)
	// if !ok || !token.Valid {
	// 	return errors.New("Invalid token")
	// }

	// now := jwt.NewNumericDate(time.Now())
	// exp, err := claims.GetExpirationTime()
	// if err != nil {
	// 	return err
	// }
	// if time.Now().After(exp.Time) {
	// 	return errors.New("token expired")
	// }

	// uid, authId, err := ExtractTokenID(c)
	// user, err := u.service.GetUserByAuth(uid, authId)

	return nil, false
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
		authId := claims["auth_uuid"].(string)
		uid := claims["sub"].(float64)
		return int64(uid), authId, nil
	}
	return 0, "", nil
}

func VerifyToken(c *gin.Context) (*jwt.Token, error) {
	tokenString := ExtractToken(c)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(config.GetConfig().GetString("AUTH_JWT_SECRET")), nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}

func ExtractTokenAuth(c *gin.Context) (*AuthDetails, error) {
	token, err := VerifyToken(c)
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(jwt.MapClaims) //the token claims should conform to MapClaims
	if ok && token.Valid {
		authUuid, ok := claims["auth_uuid"].(string) //convert the interface to string
		if !ok {
			return nil, err
		}
		userId, err := strconv.ParseUint(fmt.Sprintf("%.f", claims["user_id"]), 10, 64)
		if err != nil {
			return nil, err
		}
		return &AuthDetails{
			AuthUuid: authUuid,
			UserId:   userId,
		}, nil
	}
	return nil, err
}

func BlacklistToken(token string) error {
	_, err := config.NewRedis().Get().Do("SET", "blacklist_"+token, token, "EX", 3600*24)
	if err != nil {
		return errors.New("unable to blacklist token")
	}

	return nil
}
