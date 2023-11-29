package initializer

import "github.com/raselldev/go-jwt/models"

func SyncDatabase() {
	DB.AutoMigrate(&models.User{})
}
