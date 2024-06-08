package core

type User struct {
	Id int `db:"id"`
	// Sub is the user_id field from Auth0 and might be deprecated in the future.
	Sub  string `db:"sub"`
	Name string `db:"name"`
}
