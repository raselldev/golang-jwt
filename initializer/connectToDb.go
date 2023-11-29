package initializer

import (
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// `var DB *gorm.DB` is declaring a variable named `DB` of type `*gorm.DB`. This variable will be used
// to store the connection to the PostgreSQL database. The `*gorm.DB` type is a pointer to the
// `gorm.DB` struct, which represents a GORM database connection.
var DB *gorm.DB

// The function `ConnectToDb` connects to a PostgreSQL database using the provided environment variable
// and the GORM library.
func ConnectToDb() {
	var err error
	dsn := os.Getenv("DB")
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("Failed to connect DB")
	}
}
