package initializer

import "github.com/raselldev/go-jwt/models"

// The function `SyncDatabase` is used to automatically migrate the `User` model in the database.
func SyncDatabase() {
	DB.AutoMigrate(&models.User{})
}
