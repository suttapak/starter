package helpers

import (
	"encoding/json"

	"golang.org/x/crypto/bcrypt"
)

type (
	Helper interface {
		HashPassword(password string) (string, error)
		CheckPassword(hashPassword string, plainPassword []byte) (bool, error)
		ParseJson(src any, dns any) error
	}
	helper struct {
	}
)

// ParseJson implements Helper.
func (p helper) ParseJson(src any, dns any) error {
	b, err := json.Marshal(src)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(b, dns); err != nil {
		return err
	}
	return nil
}

// CheckPassword implements PasswordHelper.
func (p helper) CheckPassword(hashPassword string, plainPassword []byte) (bool, error) {
	hashPW := []byte(hashPassword)
	if err := bcrypt.CompareHashAndPassword(hashPW, plainPassword); err != nil {
		return false, err
	}
	return true, nil
}

// HashPassword implements PasswordHelper.
func (p helper) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 4)
	return string(bytes), err
}

func NewHelper() Helper {
	return helper{}
}
