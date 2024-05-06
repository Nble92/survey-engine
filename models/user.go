package models

import (
	"errors"

	"github.com/jinzhu/gorm"
	// the blank import is the way GORM wants you to do it
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"golang.org/x/crypto/bcrypt"
)

// User struct
type User struct {
	gorm.Model    `json:"-"`
	Username      string `gorm:"not null;unique" json:"username"`
	Password      string `gorm:"not null" json:"-"`
	IsAdmin       bool   `json:"isAdmin"`
	FullName      string `json:"fullName"`
	Email         string `json:"email"`
	Notifications bool   `json:"notifications"`
}

// UserManager manager
type UserManager struct {
	db *gorm.DB
}

// NewUserManager for handling the DB operations
func NewUserManager(db *gorm.DB) *UserManager {
	userManager := UserManager{}
	userManager.db = db

	userManager.db.AutoMigrate(&User{})

	return &userManager
}

// HasUsers simple check for existing users
func (userManager *UserManager) HasUsers() bool {
	user := User{}

	userManager.db.First(&user)

	if user.Username == "" {
		return false
	}

	return true
}

// ListUsers return all users
func (userManager *UserManager) ListUsers() (users []User) {
	err := userManager.db.Find(&users).Error

	if err != nil {
		return users
	}

	return users
}

// FindUser simple search
func (userManager *UserManager) FindUser(username string) *User {
	user := User{}

	userManager.db.Where("username=?", username).Find(&user)

	return &user
}

// RegisterUser create new user
func (userManager *UserManager) RegisterUser(username string, password string) (user *User, err error) {
	passwordHash := userManager.HashPassword(password)

	user = &User{
		Username: username,
		Password: passwordHash,
		IsAdmin:  !userManager.HasUsers(),
	}

	err = userManager.db.Create(&user).Error
	if err != nil {
		return nil, err
	}

	return user, nil
}

// HashPassword hash of the password to store
func (userManager *UserManager) HashPassword(password string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		panic("Permissions: bcrypt password hashing unsuccessful")
	}
	return string(hash)
}

// CheckPassword validates hashed password
func (userManager *UserManager) CheckPassword(hashedPassword string, password string) bool {
	if bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)) != nil {
		return false
	}
	return true
}

// UpdateUser update existing function
func (userManager *UserManager) UpdateUser(user *User) (retUser *User, err error) {
	err = userManager.db.Save(&user).Error
	if err != nil {
		return nil, err
	}

	return user, nil
}

// DeleteUser update existing function
func (userManager *UserManager) DeleteUser(user *User) (err error) {
	err = userManager.db.Delete(&user).Error
	if err != nil {
		return err
	}

	return nil
}

// UpdatePassword update password after validation
func (userManager *UserManager) UpdatePassword(user *User, existingPassword string, newPassword string) (returnUser *User, err error) {
	existingPasswordHash := userManager.HashPassword(newPassword)

	if existingPasswordHash == user.Password {
		newPasswordHash := userManager.HashPassword(newPassword)

		user.Password = newPasswordHash

		err = userManager.db.Save(&user).Error
		if err != nil {
			return nil, err
		}

		return user, nil
	}

	return nil, errors.New("existing password is not correct")
}
