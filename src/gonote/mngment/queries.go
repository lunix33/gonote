package mngment

const (
	// UserToken
	userTokenGetQuery = `
		SELECT * FROM UserToken
		WHERE UserToken.UserID = ? AND
			UserToken.Token = ?
		LIMIT 1`

	userTokenInsertQuery = `
		INSERT INTO UserToken(
			Token, Type, UserID, Expiracy, IP
		) VALUES (?, ?, ?, ?, ?)`

	userTokenRefreshQuery = `
		UPDATE UserToken SET
			Expiracy = ?,
			IP = ?
		WHERE UserToken.UserID = ? AND
			UserToken.Token = ?
		LIMIT 1`

	userTokenDeleteQuery = `
		DELETE FROM UserToken
		WHERE UserToken.UserID = ? AND
			UserToken.Token = ?
		LIMIT 1`

	// User
	userGetTokensQuery = `
		SELECT * FROM UserToken
		WHERE UserToken.UserID = ?
		ORDER BY UserToken.Expiracy`

	userGetQuery = `
		SELECT * FROM User
		WHERE User.Username = ?
		LIMIT 1`

	userAddQuery = `
		INSERT INTO User (
			ID, Username, Password, Email, IsAdmin
		) VALUES(?, ?, ?, ?, ?)`

	userDeleteQuery = `
		UPDATE User
		SET Deleted = 1
		WHERE User.ID = ?
		LIMIT 1`

	userUpdateQuery = `
		UPDATE User
		SET Username = ?,
			Password = ?,
			Email = ?,
			IsAdmin = ?
		WHERE User.ID = ?`

	// Setting
	settingGetQuery = `
		SELECT * FROM Setting
		WHERE Setting.Key = ?
		LIMIT 1`

	settingGetAllQuery = `
		SELECT * FROM Setting
		ORDER BY Setting.Key`

	settingDeleteQuery = `
		DELETE FROM Setting
		WHERE Setting.Key = ?
		LIMIT 1`

	settingUpsertQuery = `
		INSERT INTO Setting(Key, Value) VALUES (?, ?)
		ON CONFLICT(Key) DO UPDATE
		SET Value = ?`
)
