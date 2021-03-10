package teststore_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/zhainar/awesomeProject/internal/app/model"
	"github.com/zhainar/awesomeProject/internal/app/store"
	"github.com/zhainar/awesomeProject/internal/app/store/teststore"
	"testing"
)

func TestUserRepository_Create(t *testing.T) {
	s := teststore.New()

	user := model.TestUser(t)

	assert.NotNil(t, user)

	err := s.User().Create(user)

	assert.NoError(t, err)
}

func TestUserRepository_FindByEmail(t *testing.T) {
	email := "test@mail.ru"

	s := teststore.New()

	_, err := s.User().FindByEmail(email)
	assert.EqualError(t, err, store.ErrRecordNotFound.Error())

	u := model.TestUser(t)
	u.Email = email

	s.User().Create(u)

	u, err = s.User().FindByEmail(email)

	assert.NoError(t, err)
	assert.NotNil(t, u)
	assert.Equal(t, email, u.Email)
}
