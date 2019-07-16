INSERT INTO Note (ID, Title, Owner) VALUES (
    "1", "Hello World", "admin"
);

INSERT INTO NoteContent (NoteID, Version, Editor, Content) VALUES (
    "1", 1, "admin", "This is a test note."
);

INSERT INTO Tag(ID, Name) VALUES (
    "1", "Test"
);

INSERT INTO NoteTag VALUES (
    "1", "1"
);
