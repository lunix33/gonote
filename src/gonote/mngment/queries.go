package mngment

const (
	userTokenInsertQuery = `
		INSERT INTO UserToken(
			Token, Type, Owner, Expiracy, IP
		) VALUES (?, ?, ?, ?, ?)`

	userTokenRefreshQuery = `
		UPDATE UserToken SET
			Expiracy = ?,
			IP = ?
		WHERE Owner = ? AND Token = ?`

	userTokenDeleteQuery = `
		DELETE FROM UserToken
		WHERE Owner = ? AND Token = ?`
)
