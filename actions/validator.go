/*
 *
 * In The Name of God
 *
 * +===============================================
 * | Author:        Parham Alvani <parham.alvani@gmail.com>
 * |
 * | Creation Date: 02-02-2019
 * |
 * | File Name:     validator.go
 * +===============================================
 */

package actions

import (
	"gopkg.in/go-playground/validator.v9"
)

// DefaultValidator contains go-playground validator instance
type DefaultValidator struct {
	validator *validator.Validate
}

// Validate validates structs based on go-playground validator
func (cv *DefaultValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}
