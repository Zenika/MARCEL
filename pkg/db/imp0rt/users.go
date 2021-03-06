package imp0rt

import (
	"github.com/allez-chauffe/marcel/pkg/db"
	"github.com/allez-chauffe/marcel/pkg/db/users"
)

type userPassword struct {
	users.User
	PasswordHash string `json:"passwordHash"`
	PasswordSalt string `json:"passwordSalt"`
}

func Users(inputFile string) error {
	var usersPassword []userPassword

	return imp0rt(inputFile, &usersPassword, func() error {
		return importUsers(usersPassword)
	})
}

func importUsers(usersPassword []userPassword) error {
	var usersList = make([]users.User, 0, len(usersPassword))
	for _, up := range usersPassword {
		var u = up.User
		u.PasswordHash = up.PasswordHash
		u.PasswordSalt = up.PasswordSalt
		usersList = append(usersList, u)
	}
	return db.Users().UpsertAll(usersList)
}
