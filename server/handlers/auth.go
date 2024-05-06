package handlers

import (
	"errors"
	"time"

	"github.com/dchote/survey-engine/config"
	"github.com/dchote/survey-engine/models"

	"github.com/dgrijalva/jwt-go"
)

// JWTCustomClaims struct
type JWTCustomClaims struct {
	Username string `json:"username"`
	IsAdmin  bool   `json:"isAdmin"`
	jwt.StandardClaims
}

// UserToken returns a token string
func UserToken(user *models.User) (string, error) {
	// Set claims
	claims := &JWTCustomClaims{
		user.Username,
		user.IsAdmin,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(config.Config.Server.JWTSigningKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateUser validates that a user exists within the token
func ValidateUser(token *jwt.Token) (*models.User, error) {
	claims := token.Claims.(*JWTCustomClaims)

	user := userManager.FindUser(claims.Username)
	if user.Username == "" {
		return nil, errors.New("Token user not found in database")
	}

	return user, nil
}

// ValidateAdmin validates that the user is an admin
func ValidateAdmin(token *jwt.Token) (*models.User, error) {
	user, err := ValidateUser(token)

	if err != nil {
		return nil, err
	} else if !user.IsAdmin {
		return nil, errors.New("User is not an admin")
	}

	return user, nil
}
