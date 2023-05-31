CREATE TABLE IF NOT EXISTS posts(
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    kind_id uuid not null,
    CONSTRAINT fk_kind
        FOREIGN KEY(kind_id)
            REFERENCES kinds(id)
);

-- TODO: delete cascade all pics if the post has been deleted