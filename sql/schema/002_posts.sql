-- +goose Up
CREATE TABLE posts(
id UUID primary key,
created_at TIMESTAMP NOT NULL,
updated_at TIMESTAMP NOT NULL,
body TEXT NOT NULL,
user_ID UUID REFERENCES users(id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE posts;
