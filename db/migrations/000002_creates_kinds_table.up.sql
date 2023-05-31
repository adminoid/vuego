CREATE TABLE IF NOT EXISTS kinds(
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR (50) UNIQUE NOT NULL
);

-- if kind has been deleted, ignore touching posts