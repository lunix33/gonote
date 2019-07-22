-- Users
INSERT INTO "User" (
	"ID",                                   "Username", "Password", 													"Email", 			"Deleted",	"IsAdmin"
) VALUES (												/* Qwerty1! */
	"0932ae86-d5be-4652-b320-acd6b7ede7e6", "TestUser", "$2y$10$cYWUwepZUgvF7JmagtMGJ.O/ThY7P.dTM11fZHpfJJo7mqCRA9pQ6", "super@mail.com",	0,			0
);

-- Notes
INSERT INTO "Note" (
	"ID",									"Title", 			"UserID",								"Public",	"Added",							"Deleted"
) VALUES (
	"336d8617-2a68-43af-a30d-5fefbe66a8b3",	"Hello World",		"0932ae86-d5be-4652-b320-acd6b7ede7e6",	1,			"2019-07-22T13:00:00.0000-04:00",	0
), (
	"336d8617-2a68-43af-a30d-5fefbe66a8b4",	"Deleted World",	"0932ae86-d5be-4652-b320-acd6b7ede7e6",	1,			"2019-07-22T13:00:00.0000-04:00",	1
), (
	"336d8617-2a68-43af-a30d-5fefbe66a8b5",	"Own World",		"fa78fdff-bdbb-48f2-9d6c-eadeb078d4bd",	0,			"2019-07-22T13:00:00.0000-04:00",	0
);

-- NoteContents
INSERT INTO "NoteContent" (
	"NoteID",								"Version",	"Content",					"Updated"
) VALUES (
	"336d8617-2a68-43af-a30d-5fefbe66a8b3",	1,			"I welcome the",			"2019-07-22T13:00:00.0000-04:00"
), (
	"336d8617-2a68-43af-a30d-5fefbe66a8b3",	2,			"I welcome the world!",		"2019-07-22T14:00:00.0000-04:00"
), (
	"336d8617-2a68-43af-a30d-5fefbe66a8b4",	1,			"This doesn't even exists",	"2019-07-22T14:00:00.0000-04:00"
), (
	"336d8617-2a68-43af-a30d-5fefbe66a8b5",	1,			"This is my perfect note",	"2019-07-22T14:00:00.0000-04:00"
);

-- NoteTags
INSERT INTO "NoteTag" (
	"NoteID", "Name"
) VALUES (
	"336d8617-2a68-43af-a30d-5fefbe66a8b3", "#HotStuff!"
), (
	"336d8617-2a68-43af-a30d-5fefbe66a8b3", "Stuff"
), (
	"336d8617-2a68-43af-a30d-5fefbe66a8b4", "Void"
), (
	"336d8617-2a68-43af-a30d-5fefbe66a8b5", "Important"
);
