CREATE TABLE IF NOT EXISTS pics(
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    path VARCHAR (700) UNIQUE NOT NULL,
    alt VARCHAR (700),
    post_id uuid not null,
    CONSTRAINT fk_post
        FOREIGN KEY(post_id)
            REFERENCES posts(id)
);
