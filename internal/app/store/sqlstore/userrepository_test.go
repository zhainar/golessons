package sqlstore_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/zhainar/awesomeProject/internal/app/model"
	"github.com/zhainar/awesomeProject/internal/app/store"
	"github.com/zhainar/awesomeProject/internal/app/store/sqlstore"
	"testing"
)

func TestUserRepository_Create(t *testing.T) {
	db, teardown := sqlstore.TestDB(t, databaseURL)
	defer teardown("users")

	s := sqlstore.New(db)

	user := model.TestUser(t)

	assert.NotNil(t, user)

	err := s.User().Create(user)

	assert.NoError(t, err)
}

func TestUserRepository_FindByEmail(t *testing.T) {
	db, teardown := sqlstore.TestDB(t, databaseURL)
	defer teardown("users")

	email := "test@mail.ru"

	s := sqlstore.New(db)

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
	db, teardown := sqlstore.TestDB(t, databaseURL)
	defer teardown("users")

	email := "test@mail.ru"

	s := sqlstore.New(db)
	u := model.TestUser(t)
	u.Email = email

	s.User().Create(u)

	u2, err := s.User().Find(u.ID)

	assert.NoError(t, err)
	assert.NotNil(t, u2)
	assert.Equal(t, u.ID, u2.ID)
}
