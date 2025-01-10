package helpers

import (
	"encoding/json"

	"github.com/stretchr/testify/mock"
)

type (
	passowrdHelperMock struct {
		mock.Mock
	}
)

// ParseJson implements Helper.
func (p *passowrdHelperMock) ParseJson(src any, dns any) error {
	args := p.Called()

	b, err := json.Marshal(src)
	if err != nil {
		return args.Error(0)
	}
	if err := json.Unmarshal(b, dns); err != nil {
		return args.Error(0)
	}
	return args.Error(0)
}

// CheckPassword implements PasswordHelper.
func (p *passowrdHelperMock) CheckPassword(hashPassword string, plainPassword []byte) (bool, error) {
	args := p.Called()
	return args.Bool(0), args.Error(1)
}

// HashPassword implements PasswordHelper.
func (p *passowrdHelperMock) HashPassword(password string) (string, error) {
	args := p.Called()
	return args.String(0), args.Error(1)
}

func NewHelperMock() *passowrdHelperMock {
	return &passowrdHelperMock{}
}
