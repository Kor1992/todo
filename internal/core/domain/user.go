package domain

import (
	"fmt"
	"regexp"

	core_errors "github.com/Kor1992/todo/internal/core/errors"
)

type User struct {
	ID      int
	Version int

	FullName    string
	PhoneNumber *string
}

func NewUserUnInitialized(fullName string, phoneNumber *string) User {
	return NewUser(fullName, phoneNumber, UninitializedID, UninitializedVersion)
}

func NewUser(fullName string, phoneNumber *string, id int, version int) User {
	return User{
		FullName:    fullName,
		PhoneNumber: phoneNumber,
		ID:          id,
		Version:     version,
	}
}

func (u *User) Validate() error {
	fullNameLenght := len([]rune(u.FullName))
	if fullNameLenght < 3 || fullNameLenght > 100 {
		return fmt.Errorf("invalid fullName len:%d: %w", fullNameLenght, core_errors.ErrInvalidArgument)
	}

	if u.PhoneNumber != nil {
		phoneNumberLen := len([]rune(*u.PhoneNumber))
		if phoneNumberLen < 10 || phoneNumberLen > 15 {
			return fmt.Errorf("invalid phone number len:%d: %w", phoneNumberLen, core_errors.ErrInvalidArgument)

		}
		re := regexp.MustCompile(`^\+?[0-9]+$`)
		if !re.MatchString(*u.PhoneNumber) {
			return fmt.Errorf("invalid phone number format:%w", core_errors.ErrInvalidArgument)
		}
	}
	return nil
}
