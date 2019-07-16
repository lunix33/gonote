package user

// User represent a user of the platform.
type User struct {
	Username string
	Password string
	Email    string
	Deleted  bool
}
