package controllers

import (
	"fmt"
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
	var body struct {
		Email    string
		Password string
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to read body",
		})
		return
	}

	validEmail := helpers.IsEmailValid(body.Email)

	if !validEmail {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "email format is invalid",
		})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to hash password",
		})
		return
	}

	user := models.User{Email: body.Email, Password: string(hash)}
	result := initializer.DB.Create(&user)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to create user",
		})

		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "OK",
	})
}

// The Login function handles user authentication by checking the email and password, generating a JWT
// token, and returning it to the user.
func Login(c *gin.Context) {
	var body struct {
		Email    string
		Password string
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to read body",
		})

		return
	}

	var user models.User
	initializer.DB.First(&user, "email=?", body.Email)

	if user.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid email or password",
		})
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid email or password",
		})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	})

	secret := []byte(os.Getenv("SECRET"))
	tokenString, err := token.SignedString(secret)

	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to create token",
		})
		return
	}

	user.Token = tokenString
	user.TokenExpire = time.Now().Add(time.Hour * 24 * 30).Unix()
	initializer.DB.Save(&user)

	c.JSON(http.StatusOK, gin.H{
		"token":   tokenString,
		"expires": time.Now().Add(time.Hour * 24 * 30).Unix(),
	})
}
