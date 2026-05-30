-- name: CreatePost :one
INSERT INTO posts (id, created_at, updated_at, body, user_id)
VALUES (
	gen_random_uuid(), NOW(), NOW(), $1, $2
	)
	RETURNING *;

-- name: ResetPosts :exec
DELETE FROM posts;


-- name: GetAllPosts :many
Select * FROM posts 
ORDER BY created_at ASC;
