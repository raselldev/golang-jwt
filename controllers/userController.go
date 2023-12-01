package controllers

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/raselldev/go-jwt/helpers"
	"github.com/raselldev/go-jwt/initializer"
	"github.com/raselldev/go-jwt/models"
	"golang.org/x/crypto/bcrypt"
)

// The SignUp function handles the sign-up process by reading the request body, hashing the password,
// creating a new user in the database, and returning a response.
func SignUp(c *gin.Context) {
	var requestBody struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.Bind(&requestBody); err != nil {
		handleError(c, http.StatusBadRequest, "failed to read body")
		return
	}

	if !helpers.IsEmailValid(requestBody.Email) {
		handleError(c, http.StatusBadRequest, "email format is invalid")
		return
	}

	hashedPassword, err := hashPassword(requestBody.Password)
	if err != nil {
		handleError(c, http.StatusBadRequest, "failed to hash password")
		return
	}

	user := models.User{Email: requestBody.Email, Password: hashedPassword}
	if err := createNewUser(&user); err != nil {
		handleError(c, http.StatusBadRequest, "failed to create user")
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "OK",
	})
}

func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(hashedPassword), err
}

func createNewUser(user *models.User) error {
	result := initializer.DB.Create(user)
	return result.Error
}

// The Login function handles user authentication by checking the email and password, generating a JWT
// token, and returning it to the user.
func Login(c *gin.Context) {
	var requestBody struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.Bind(&requestBody); err != nil {
		handleError(c, http.StatusBadRequest, "failed to read body")
		return
	}

	user, err := findUserByEmail(requestBody.Email)
	if err != nil || user.ID == 0 || !checkPassword(user.Password, requestBody.Password) {
		handleError(c, http.StatusBadRequest, "invalid email or password")
		return
	}

	tokenString, err := generateToken(user.ID)
	if err != nil {
		handleError(c, http.StatusBadRequest, "failed to create token")
		return
	}

	updateUserToken(user, tokenString)

	c.JSON(http.StatusOK, gin.H{
		"token":   tokenString,
		"expires": user.TokenExpire,
	})
}

func handleError(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, gin.H{"error": message})
}

func findUserByEmail(email string) (models.User, error) {
	var user models.User
	err := initializer.DB.First(&user, "email=?", email).Error
	return user, err
}

func checkPassword(hash, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}

func generateToken(userID uint) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	})

	secret := []byte(os.Getenv("SECRET"))
	return token.SignedString(secret)
}

func updateUserToken(user models.User, tokenString string) {
	user.Token = tokenString
	user.TokenExpire = time.Now().Add(time.Hour * 24 * 30).Unix()
	initializer.DB.Save(&user)
}
