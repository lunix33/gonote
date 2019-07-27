BEGIN TRANSACTION;

-- Create application tables.
CREATE TABLE "User" (
	"ID"		TEXT	NOT NULL,
	"Username"	TEXT	NOT NULL,
	"Password"	TEXT	NOT NULL,
	"Email"		TEXT	NOT NULL	DEFAULT "no-mail@domain",
	"Deleted"	BOOLEAN	NOT NULL	DEFAULT 0,
	"IsAdmin"	BOOLEAN NOT NULL	DEFAULT 0,
	PRIMARY KEY("ID"),
	UNIQUE("Username")
);

CREATE TABLE "UserToken" (
	"Token"		TEXT		NOT NULL,
	"Type"		TEXT		NOT NULL,
	"UserID"	TEXT		NOT NULL,
	"Created"	DATETIME	NOT NULL	DEFAULT (datetime()),
	"Expiry"	DATETIME	NOT NULL,
	"IP"		TEXT		NOT NULL,
	PRIMARY KEY("Token"),
	FOREIGN KEY("UserID") REFERENCES "User"("ID") ON DELETE CASCADE 
);

CREATE TABLE "Note" (
	"ID"		TEXT		NOT NULL,
	"Title"		TEXT		NOT NULL,
	"UserID"	TEXT		NOT NULL,
	"Public"	BOOLEAN		NOT NULL	DEFAULT 0,
	"Added"		DATETIME				DEFAULT (datetime()),
	"Deleted"	INTEGER					DEFAULT 0,
	PRIMARY KEY("ID"),
	FOREIGN KEY("UserID") REFERENCES "User"("ID") ON DELETE CASCADE
);

CREATE TABLE "NoteContent" (
	"NoteID"	TEXT		NOT NULL,
	"Version"	INTEGER		NOT NULL,
	"Content"	TEXT		NOT NULL,
	"Updated"	DATETIME	NOT NULL	DEFAULT (datetime()),
	PRIMARY KEY("NoteID", "Version"),
	FOREIGN KEY("NoteID") REFERENCES "Note"("ID") ON DELETE CASCADE
);


CREATE TABLE "NoteTag" (
	"NoteID"	TEXT	NOT NULL,
	"Name"		TEXT	NOT NULL,
	PRIMARY KEY("NoteID", "Name"),
	FOREIGN KEY("NoteID") REFERENCES "Note"("ID")
);

CREATE TABLE "Setting" (
	"Key"	TEXT	NOT NULL,
	"Value"	TEXT	NOT NULL,
	PRIMARY KEY("Key")
);

-- Create a default user.
INSERT INTO "User" (
	"ID",									"Username",	"Password",																"Email",				"IsAdmin"
) VALUES(												/* alpine */
	"fa78fdff-bdbb-48f2-9d6c-eadeb078d4bd",	"admin",	"$2y$10$.3wZ.zH0A5WGfQRhI59uOeQo2uhzqxc5.qXkQcS6kYB1C4OGjMiX6",			"no-mail@domain.com",	1
);

-- Set default application settings
INSERT INTO "Setting"(
		"Key", 			"Value"
) VALUES
	( 	"DBVersion", 	"1" ),
	( 	"CustomPath", 	"custom" ),
	( 	"Port", 		"8080" ),
	( 	"Interface", 	"localhost" ),
	(	"SiteTitle",	"goNote" );

COMMIT;
