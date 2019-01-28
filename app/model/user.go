package model

import "extension/database"

type User struct {
	ID   int
	Name string
	PWD  string
}

func (this *User) Orm() *database.Sql {

	return database.Table("user")
}
