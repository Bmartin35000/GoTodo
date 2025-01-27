CREATE TABLE todos
(
    id uuid PRIMARY KEY NOT NULL,
    title VARCHAR(100),
    completed BOOL,
    createdat DATE
)
