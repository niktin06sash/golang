CREATE TABLE Userz (
    Userid UUID PRIMARY KEY,
    Username text,
    UserEmail text unique,
    UserPassword text
)