CREATE TABLE Userz (
    Userid uuid.UUID PRIMARY KEY,
    Username text,
    UserEmail text unique,
    UserPassword text
)