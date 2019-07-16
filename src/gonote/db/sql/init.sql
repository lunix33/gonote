BEGIN TRANSACTION;

-- Create application tables.
CREATE TABLE User (
	Username	TEXT	NOT NULL,
	Password	TEXT	NOT NULL,
	Email		TEXT,
	Deleted		INTEGER	NOT NULL	DEFAULT 0,
	IsAdmin		INTEGER NOT NULL	DEFAULT 0,
	PRIMARY KEY(Username)
);

CREATE TABLE Note (
	ID			TEXT		NOT NULL,
	Title		TEXT		NOT NULL,
	Owner		TEXT		NOT NULL,
	Added		DATETIME				DEFAULT (datetime()),
	Deleted		INTEGER					DEFAULT 0,
	PRIMARY KEY(ID),
	FOREIGN KEY(Owner) REFERENCES User(Username) ON DELETE CASCADE
);

CREATE TABLE NoteContent (
	NoteID	TEXT		NOT NULL,
	Version	INTEGER		NOT NULL,
	Editor	TEXT		NOT NULL,
	Content	TEXT		NOT NULL,
	Updated	DATETIME	NOT NULL	DEFAULT (datetime()),
	PRIMARY KEY(NoteID, Version),
	FOREIGN KEY(Editor) REFERENCES User(Username) ON DELETE CASCADE, 
	FOREIGN KEY(NoteID) REFERENCES Note(ID) ON DELETE CASCADE
);


CREATE TABLE NoteTag (
	NoteID	TEXT	NOT NULL,
	Name	TEXT	NOT NULL,
	PRIMARY KEY(NoteID, Name),
	FOREIGN KEY(NoteID) REFERENCES Note(ID)
);

CREATE TABLE Setting (
	Key		TEXT	NOT NULL,
	Value	TEXT	NOT NULL,
	PRIMARY KEY(Key)
);

-- Create a default user.
INSERT INTO User(Username, Password, Email, IsAdmin) VALUES(
	"admin",
	"",
	"no-mail@domain.com",
	1
);

-- Set default application settings
INSERT INTO Setting(Key, Value) VALUES
	( "DBVersion", "1" ),
	( "CustomPath", "custom" ),
	( "Port", "8080" ),
	( "Interface", "localhost" );

COMMIT;
