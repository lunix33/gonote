package mngment

const (
	// ---- UserToken ---------------------------------------------------------
	// Query to get a specific user's token.
	// Params: User ID, Token
	userTokenGetQuery = `
		SELECT UserToken.*
		FROM UserToken
		WHERE UserToken.UserID = ? AND
			UserToken.Token = ?
		LIMIT 1`

	// Query to add a new token.
	// Params: Token, Token type, User ID, Expiracy, User's IP
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

	// ---- User --------------------------------------------------------------
	userGetTokensQuery = `
		SELECT UserToken.*
		FROM UserToken
		WHERE UserToken.UserID = ?
		ORDER BY UserToken.Expiracy`

	userGetQuery = `
		SELECT User.*
		FROM User
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

	// ---- Setting -----------------------------------------------------------
	settingGetQuery = `
		SELECT Setting.*
		FROM Setting
		WHERE Setting.Key = ?
		LIMIT 1`

	settingGetAllQuery = `
		SELECT Setting.*
		FROM Setting
		ORDER BY Setting.Key`

	settingDeleteQuery = `
		DELETE FROM Setting
		WHERE Setting.Key = ?
		LIMIT 1`

	settingUpsertQuery = `
		INSERT INTO Setting(Key, Value) VALUES (?, ?)
		ON CONFLICT(Key) DO UPDATE
		SET Value = ?`

	// ---- Tag ---------------------------------------------------------------
	tagGetAllQuery = `
		SELECT DISTINCT Name
		FROM NoteTag
		ORDER BY NoteTag.Name`

	tagAddQuery = `
		INSERT INTO NoteTag(
			NoteID, Name
		) VALUES (?, ?)`

	tagGetNotesQuery = `
		SELECT Note.*
		FROM Note
		INNER JOIN (
			SELECT NoteContent.Updated, NoteContent.NoteID
			FROM NoteContent
			ORDER BY NoteContent.Version DESC
			LIMIT 1
		) nc ON nc.NoteID = Note.ID
		WHERE Note.ID = ?
		ORDER BY nc.Updated DESC`

	tagRemoveQuery = `
		DELETE FROM NoteTag
		WHERE NoteTag.NoteID = ? AND
			NoteTag.Name = ?`

	// ---- Note --------------------------------------------------------------
	noteGetQuery = ``

	noteGetAllQuery = ``

	noteAddQuery = ``

	// Query used to delete a note.
	// Params : Note.ID
	noteDeleteQuery = `
		UPDATE Note
		SET Deleted = 1
		WHERE Note.ID = ?`

	noteUpdateQuery = ``

	noteGetAllTagsQuery = ``

	// ---- NoteContent -------------------------------------------------------
)
