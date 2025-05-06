package service

// Instructions:
// Do not manually register services here unless absolutely necessary.
// This code is designed to auto-register services dynamically.
// Any manual changes might interfere with the auto-registration process.
// Ensure that any new service follows the auto-registration conventions.

import (
	"gorm.io/gorm"
)

type Container struct {
}

func NewContainer(db *gorm.DB) *Container {
	return &Container{
		// register auto
	}
}
