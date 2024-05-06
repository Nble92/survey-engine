package handlers

import (
	//"encoding/json"
	//"fmt"
	//"io/ioutil"
	"log"
	"net/http"

	"github.com/dchote/survey-engine/models"

	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
)

var (
	userManager *models.UserManager
)

// AuthRequest struct
type AuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// AuthResponse struct
type AuthResponse struct {
	Token    string `json:"token"`
	Username string `json:"Username"`
	IsAdmin  bool   `json:"isAdmin"`
}

// UserObject struct
type UserObject struct {
	Username      string `json:"username"`
	IsAdmin       bool   `json:"isAdmin"`
	FullName      string `json:"fullName"`
	Email         string `json:"email"`
	Notifications bool   `json:"notifications"`
}

// UserObjectList struct
type UserObjectList struct {
	Users []UserObject `json:"users"`
}

// UserUpdateRequest struct
type UserUpdateRequest struct {
	IsAdmin       bool   `json:"isAdmin"`
	FullName      string `json:"fullName"`
	Email         string `json:"email"`
	Notifications bool   `json:"notifications"`
}

// InitUserHandler initialize the user manager
func InitUserHandler(db *gorm.DB) {
	userManager = models.NewUserManager(db)
}

// Login login method
func Login() echo.HandlerFunc {
	return func(c echo.Context) (err error) {
		loginUser := new(AuthRequest)

		if err = c.Bind(loginUser); err != nil {
			return
		}

		if loginUser.Username == "" || loginUser.Password == "" {
			return echo.NewHTTPError(http.StatusUnauthorized, "Please provide valid credentials")
		}

		user := userManager.FindUser(loginUser.Username)

		if user.Username == "" {
			log.Printf("Username '%s' not found", loginUser.Username)
			return echo.NewHTTPError(http.StatusUnauthorized, "Please provide valid credentials")
		}

		if !userManager.CheckPassword(user.Password, loginUser.Password) {
			log.Printf("Password for '%s' was not correct", loginUser.Username)
			return echo.NewHTTPError(http.StatusUnauthorized, "Please provide valid credentials")
		}

		// generate token
		token, err := UserToken(user)

		if err != nil {
			log.Printf("Error generating token for user '%s': %s", loginUser.Username, err)
			return echo.NewHTTPError(http.StatusInternalServerError, "Error generating token")
		}

		response := &AuthResponse{
			Token:    token,
			Username: user.Username,
			IsAdmin:  user.IsAdmin,
		}

		return c.JSON(http.StatusOK, response)
	}
}

// Register register method
func Register() echo.HandlerFunc {
	return func(c echo.Context) (err error) {
		registerUser := new(AuthRequest)

		if err = c.Bind(registerUser); err != nil {
			return
		}

		if registerUser.Username == "" || registerUser.Password == "" {
			return echo.NewHTTPError(http.StatusUnauthorized, "Please provide valid credentials")
		}

		user, err := userManager.RegisterUser(registerUser.Username, registerUser.Password)

		if err != nil {
			log.Printf("Error creating user: %s", err)
			return echo.NewHTTPError(http.StatusUnauthorized, "Error creating user, please ensure you have a unique username")
		}

		// generate token
		token, err := UserToken(user)

		if err != nil {
			log.Printf("Error generating token for user '%s': %s", registerUser.Username, err)
			return echo.NewHTTPError(http.StatusInternalServerError, "Error generating token")
		}

		response := &AuthResponse{
			Token:    token,
			Username: user.Username,
			IsAdmin:  user.IsAdmin,
		}

		return c.JSON(http.StatusOK, response)
	}
}

// UserInfo fetch user info
func UserInfo() echo.HandlerFunc {
	return func(c echo.Context) (err error) {
		user, err := ValidateUser(c.Get("user").(*jwt.Token))

		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, "Invalid Token")
		}

		response := &AuthResponse{
			Username: user.Username,
			IsAdmin:  user.IsAdmin,
		}

		return c.JSON(http.StatusOK, response)
	}
}

// ListUsers list all users
func ListUsers() echo.HandlerFunc {
	return func(c echo.Context) error {
		_, err := ValidateAdmin(c.Get("user").(*jwt.Token))
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, "Permission Denied")
		}

		users := userManager.ListUsers()

		userObjectList := new(UserObjectList)

		for _, user := range users {
			userObject := UserObject{
				Username: user.Username,
				IsAdmin:  user.IsAdmin,
			}
			userObjectList.Users = append(userObjectList.Users, userObject)
		}

		return c.JSON(http.StatusOK, userObjectList)
	}
}

// FetchUser list all users
func FetchUser() echo.HandlerFunc {
	return func(c echo.Context) error {
		_, err := ValidateAdmin(c.Get("user").(*jwt.Token))
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, "Permission Denied")
		}

		username := c.Param("name")

		user := userManager.FindUser(username)

		if user.Username == "" {
			return echo.NewHTTPError(http.StatusNotFound, "User not found")
		}

		userObject := UserObject{
			Username:      user.Username,
			IsAdmin:       user.IsAdmin,
			FullName:      user.FullName,
			Email:         user.Email,
			Notifications: user.Notifications,
		}

		return c.JSON(http.StatusOK, userObject)
	}
}

// UpdateUser update function
func UpdateUser() echo.HandlerFunc {
	return func(c echo.Context) error {
		_, err := ValidateAdmin(c.Get("user").(*jwt.Token))
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, "Permission Denied")
		}

		username := c.Param("name")

		userUpdateRequest := new(UserUpdateRequest)
		if err = c.Bind(userUpdateRequest); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}

		user := userManager.FindUser(username)

		if user.Username == "" {
			return echo.NewHTTPError(http.StatusNotFound, "User not found")
		}

		user.IsAdmin = userUpdateRequest.IsAdmin
		user.FullName = userUpdateRequest.FullName
		user.Email = userUpdateRequest.Email
		user.Notifications = userUpdateRequest.Notifications

		userManager.UpdateUser(user)

		return c.JSON(http.StatusOK, false)
	}
}

// DeleteUser delete function
func DeleteUser() echo.HandlerFunc {
	return func(c echo.Context) error {
		_, err := ValidateAdmin(c.Get("user").(*jwt.Token))
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, "Permission Denied")
		}

		username := c.Param("name")

		user := userManager.FindUser(username)

		if user.Username == "" {
			return echo.NewHTTPError(http.StatusNotFound, err)
		}

		userManager.DeleteUser(user)

		return c.JSON(http.StatusOK, true)
	}
}
