package mngment

import (
	"gonote/db"
	"log"

	"golang.org/x/crypto/bcrypt"
)

/*
	Security procedure for users:
	First to authenticate the user a call must be placed to `POST /api/login`
	The server will reply with a token.

	In the subsquant requests the user must include the token in the header as follow:
	Authorization: Token base64("username:token")
*/

// SetPassword takes the plain password of the user and encrypt it before it is stored.
// `p`: Is the user's password in plain-text.
func (u *User) SetPassword(p string) {
	pba := []byte(p)
	c, cerr := bcrypt.GenerateFromPassword(pba, bcrypt.DefaultCost)
	if cerr != nil {
		log.Fatalln("Unable to encrypt the password.")
		return
	}
	u.Password = string(c)
}

// ComparePassword compare an encrypted password with a pain-text password.
// `p`: is the user's password in plain-text.
// Returns true if the password are a match, otherwise false.
func (u *User) ComparePassword(p string) bool {
	var (
		pba = []byte(p)
		cba = []byte(u.Password)
		err error
	)
	err = bcrypt.CompareHashAndPassword(cba, pba)
	return err == nil
}

func (u *User) GetTokens(c *db.Conn) []*UserToken {

}
