CREATE TABLE IF NOT EXISTS users(
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR (50) UNIQUE NOT NULL,
    password VARCHAR (50) NOT NULL,
    email VARCHAR (300) UNIQUE NOT NULL
);

INSERT INTO users (name, password, email)
VALUES ('ivan', 'Народ', 'iv@mail.ru');

INSERT INTO users (name, password, email)
VALUES ('petr', 'пароль', 'p@mail.ru');

INSERT INTO users (name, password, email)
VALUES ('adminoid', 'passwork', 'adminoid@mail.ru');

