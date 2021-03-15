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

func TestUserRepository_Find(t *testing.T) {
	s := teststore.New()
	u := model.TestUser(t)

	s.User().Create(u)

	u2, err := s.User().Find(u.ID)

	assert.NoError(t, err)
	assert.NotNil(t, u2)
	assert.Equal(t, u.ID, u2.ID)
}
