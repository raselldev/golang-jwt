package models

import (
	"gorm.io/gorm"
)

// The User struct represents a user with email, password, and token fields.
// @property  - - `gorm.Model`: This is a struct provided by the GORM library that includes common
// fields like `ID`, `CreatedAt`, `UpdatedAt`, and `DeletedAt`. It is used to automatically add these
// fields to the `User` struct.
// @property {string} Email - The `Email` property is a string that represents the email address of the
// user. It is tagged with `gorm:"unique"` to ensure that each email address is unique in the database.
// @property {string} Password - The "Password" property is a string that represents the user's
// password. It is used for authentication and should be securely stored and encrypted.
// @property {string} Token - The "Token" property is a string that is used for authentication and
// authorization purposes. It is typically a randomly generated string that is unique to each user.
// This token can be used to verify the identity of the user and grant access to certain resources or
// actions.
type User struct {
	gorm.Model
	Email       string `gorm:"unique"`
	Password    string
	Token       string
	TokenExpire int64
}
