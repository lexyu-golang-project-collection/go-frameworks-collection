-- name: GetAuthor :one
SELECT *
FROM authors
WHERE id = ?;

-- name: ListAuthors :many
SELECT *
FROM authors
ORDER BY name;

-- name: CreateAuthor :execresult
INSERT INTO authors (name, bio)
VALUES (?, ?);

-- name: UpdateAuthor :exec
UPDATE authors
set name = ?,
    bio = ?
WHERE id = ?;

-- name: DeleteAuthor :exec
DELETE FROM authors
WHERE id = ?;

-- name: GetAuthorById :one
SELECT * FROM authors
WHERE id = LAST_INSERT_ID();

-- name: GetUpdatedAuthor :one
SELECT * FROM authors
WHERE id = ?;