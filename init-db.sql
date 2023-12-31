CREATE TABLE Users (
    ID int IDENTITY(1,1),
    Email varchar(255),
    FirstName varchar(255),
	LastName varchar(255),
	Password varchar(255),
    Active bit,
    CreatedAt DateTime2,
	UpdatedAt DateTime2
);

INSERT INTO Users
VALUES ('admin@example.com', 'Admin', 'User', '$2a$12$1zGLuYDDNvATh4RA4avbKuheAMpb1svexSzrQm7up.bnpwQHs0jNe', 1, GETDATE(), GETDATE())

-- password is 'verysecret'