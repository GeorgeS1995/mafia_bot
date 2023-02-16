package test_db

import (
	"github.com/GeorgeS1995/mafia_bot/internal/db"
	"github.com/GeorgeS1995/mafia_bot/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

type DBPackageTestSuite struct {
	test.DBTestSuite
}

func TestDBPackage(t *testing.T) {
	suite.Run(t, new(DBPackageTestSuite))
}

func (suite *DBPackageTestSuite) TestUpdateOrCreateUserExistingUserOK() {
	defer suite.Tx.Rollback()
	polemicaID := test.RandStringRunes(3)
	polemicaNickName := test.RandStringRunes(3)
	polemicaNickNameNew := test.RandStringRunes(3)
	existingUser := &db.User{
		PolemicaId:       polemicaID,
		PolemicaNickName: polemicaNickName,
	}
	suite.Tx.Create(existingUser)
	creatAt := existingUser.CreatedAt

	user, err := db.UpdateOrCreateUser(suite.Tx, &db.User{
		PolemicaId:       polemicaID,
		PolemicaNickName: polemicaNickNameNew,
	})

	timeNow := time.Now()
	if err != nil {
		suite.T().Fatal("Update user error, ", err)
	}
	assertedUser := &db.User{}
	var count int64
	suite.Tx.Model(&db.User{}).Where("polemica_id = ?", polemicaID).Count(&count).First(assertedUser)
	assert.Equal(suite.T(), int(count), 1)
	assert.Equal(suite.T(), user.PolemicaId, assertedUser.PolemicaId)
	assert.Equal(suite.T(), assertedUser.PolemicaNickName, polemicaNickNameNew)
	assert.Equal(suite.T(), assertedUser.CreatedAt.Round(time.Second), creatAt.Round(time.Second))
	assert.Equal(suite.T(), assertedUser.UpdatedAt.Round(time.Minute), timeNow.Round(time.Minute))
}

func (suite *DBPackageTestSuite) TestUpdateOrCreateUserCreateOK() {
	defer suite.Tx.Rollback()
	polemicaID := test.RandStringRunes(3)
	polemicaNickName := test.RandStringRunes(3)

	user, err := db.UpdateOrCreateUser(suite.Tx, &db.User{
		PolemicaId:       polemicaID,
		PolemicaNickName: polemicaNickName,
	})

	timeNow := time.Now()
	if err != nil {
		suite.T().Fatal("Create user error, ", err)
	}
	assertedUser := &db.User{}
	var count int64
	suite.Tx.Model(&db.User{}).Where("polemica_id = ?", polemicaID).Count(&count).First(assertedUser)
	assert.Equal(suite.T(), int(count), 1)
	assert.Equal(suite.T(), user.PolemicaId, assertedUser.PolemicaId)
	assert.Equal(suite.T(), assertedUser.PolemicaNickName, polemicaNickName)
	assert.Equal(suite.T(), assertedUser.UpdatedAt.Round(time.Minute), timeNow.Round(time.Minute))
	assert.Equal(suite.T(), assertedUser.CreatedAt.Round(time.Minute), timeNow.Round(time.Minute))
}
