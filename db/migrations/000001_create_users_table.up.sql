CREATE TABLE IF NOT EXISTS users(
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR (50) UNIQUE NOT NULL,
    email VARCHAR (300) UNIQUE NOT NULL,
    password_hash bytea,
    refresh_token VARCHAR (300)
);

-- INSERT INTO users (name, email)
-- VALUES ('ivan', 'iv@mail.ru');
--
-- INSERT INTO users (name, email)
-- VALUES ('petr', 'p@mail.ru');
--
-- INSERT INTO users (name, email)
-- VALUES ('adminoid', 'adminoid@mail.ru');

