package services

import (
	"bigmind/xcheck-be/config"
	"bigmind/xcheck-be/internal/dto"
	"bigmind/xcheck-be/internal/models"
	"bigmind/xcheck-be/internal/repositories"
	"bigmind/xcheck-be/utils"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/datatypes"
)

type AuthService struct {
	u repositories.UserRepository
}

// var roleInstance *RoleService

// NewAuthService constructs a new AuthService using the provided UserRepository.
//
// Parameter:
//
//	u - UserRepository: the user repository to be injected into the AuthService.
//
// Return type:
//
//	*AuthService: the newly created AuthService instance.
func NewAuthService(u repositories.UserRepository) *AuthService {
	return &AuthService{
		u: u,
	}
}

func (s *AuthService) GetUserByUsername(uname string) (models.User, error) {
	result, err := s.u.FindByUsername(uname)
	if err == nil {
		result.Password = ""
		result.AuthUuids = nil
	}
	return result, err
}

func (s *AuthService) GetUserByID(id int64) (models.User, error) {
	result, err := s.u.FindByID(id)
	if err == nil {
		result.Password = ""
		result.AuthUuids = nil
	}
	return result, err
}

func (s *AuthService) GetUserByAuth(id int64, authID string) (models.User, error) {
	//result, err := s.u.FindByAuth(id, authID)
	result, err := s.u.FindByID(id)
	if err == nil {
		result.Password = ""
		result.AuthUuids = nil
	}
	return result, err
}

func (s *AuthService) Signin(username string, password string) (dto.UserLoginResponse, map[string]string, error) {
	user, err := s.u.FindByUsername(username)
	if err != nil {
		return dto.UserLoginResponse{}, nil, errors.New("invalid username or password")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return dto.UserLoginResponse{}, nil, errors.New("invalid username or password")
	}
	// create auth uuid
	authD, _ := s.u.CreateAuth(user.ID)

	tokenPair, err := utils.GenerateTokenUUID(&authD)
	if err != nil {
		return dto.UserLoginResponse{}, nil, err
	}
	user.Password = ""
	user.AuthUuids = nil

	menuGroup := utils.RoleID(*user.RoleID).GetMenu()
	jsonData, err := json.Marshal(menuGroup)
	if err != nil {
		fmt.Println(err)
		return dto.UserLoginResponse{}, nil, err
	}

	response := dto.UserLoginResponse{
		ID:        user.ID,
		Username:  user.Username,
		Password:  user.Password,
		Email:     user.Email,
		RoleID:    user.RoleID,
		AuthUuids: user.AuthUuids,
		MenuGroup: datatypes.JSON(jsonData),
	}

	return response, tokenPair, nil
}

func (s *AuthService) RefreshToken(refreshToken string) (map[string]string, error) {
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(config.GetConfig().GetString("AUTH_JWT_SECRET")), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		fmt.Println(claims)
		// Get the user record from database or
		// run through your business logic to verify if the user can log in
		avail, err := s.u.FindByID(int64(claims["sub"].(float64)))
		if err != nil {
			return nil, err
		}
		newTokenPair, err := utils.GenerateToken(&avail)
		if err != nil {
			return nil, err
		}
		return newTokenPair, nil
	}

	return nil, err
}

func (s *AuthService) ResetPassword(username string, password string) error {
	return nil
}
