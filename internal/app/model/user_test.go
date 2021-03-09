package model_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/zhainar/awesomeProject/internal/app/model"
	"testing"
)

func TestUser_BeforeCreate(t *testing.T) {
	u := model.TestUser(t)
	assert.NoError(t, u.BeforeCreate())
	assert.NotEmpty(t, u.EncryptedPassword)
}

func TestUser_Validate(t *testing.T) {
	testCases := []struct {
		name    string
		u       func() *model.User
		isValid bool
	}{
		{
			name: "valid",
			u: func() *model.User {
				return model.TestUser(t)
			},
			isValid: true,
		},
		{
			name: "with encrypted password",
			u: func() *model.User {
				user := model.TestUser(t)
				user.Password = ""
				user.EncryptedPassword = "encrypted_string"

				return user
			},
			isValid: true,
		},
		{
			name: "empty email",
			u: func() *model.User {
				user := model.TestUser(t)
				user.Email = ""
				return user
			},
			isValid: false,
		},
		{
			name: "invalid email",
			u: func() *model.User {
				user := model.TestUser(t)
				user.Email = "invalid"
				return user
			},
			isValid: false,
		},
		{
			name: "empty password",
			u: func() *model.User {
				user := model.TestUser(t)
				user.Password = ""
				return user
			},
			isValid: false,
		},
		{
			name: "short password",
			u: func() *model.User {
				user := model.TestUser(t)
				user.Password = "123"
				return user
			},
			isValid: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.isValid {
				assert.NoError(t, tc.u().Validate())
			} else {
				assert.Error(t, tc.u().Validate())
			}
		})
	}

	u := model.TestUser(t)
	assert.NoError(t, u.Validate())
}
