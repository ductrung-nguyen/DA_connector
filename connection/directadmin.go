package connection

const ACCOUNT_TYPE_ADMIN = "admin"
const ACCOUNT_TYPE_RESELLER = "reseller"
const ACCOUNT_TYPE_USER = "user"

type DirectAdmin struct {
	AuthenticatedUser string

	Username string

	Password string

	BaseURL string

	Connection string
}

// func (d *DirectAdmin) AddLoginKey(keyName string, keyValue string) {
// 	query := map[string]string{
// 		"action":  "create",
// 		"keyname": keyName,
// 	}
// }
