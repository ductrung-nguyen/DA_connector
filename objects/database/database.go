package database

import (
	"da_connector/context"
	"da_connector/objects/base"
	"da_connector/objects/users"
)

const CACHE_ACCESS_HOSTS = "access_hosts"

type Database struct {
	base.BaseObject
	owner        users.User
	databaseName string
}

func New(name string, owner users.User, context context.UserContext) Database {
	db := Database{
		BaseObject: base.New(name, context),
		owner:      owner,
	}
	db.databaseName = owner.GetUsername() + "_" + db.GetName()
	return db
}

// CreateDBForUser Creates a new database under the specified user.
func CreateDBForUser(user users.User, name string, username string, password string) {
	payload := map[string]string{
		"action": "create",
		"name":   name,
	}
	if len(password) > 0 {
		payload["password"] = password
	} else {
		payload["userlist"] = username
	}

	// user.Context
}
