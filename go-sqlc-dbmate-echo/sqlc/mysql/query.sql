-- name: CreateAuthor :execresult
INSERT INTO authors (
    name, bio
) VALUES (
    ?, ?
);

-- name: GetAuthor :one
SELECT * FROM authors
WHERE id = ? LIMIT 1;

-- name: ListAuthors :many
SELECT * FROM authors
ORDER BY name;

-- name: UpdateAuthor :execresult
UPDATE authors
SET name = ?,
    bio = ?
WHERE id = ?;

-- name: DeleteAuthor :exec
DELETE FROM authors
WHERE id = ?;

-- name: CreateBook :execresult
INSERT INTO books (
    title, author_id, published_year, genre, description
) VALUES (
    ?, ?, ?, ?, ?
);

-- name: GetBook :one
SELECT b.*, a.name as author_name
FROM books b
JOIN authors a ON b.author_id = a.id
WHERE b.id = ? LIMIT 1;

-- name: ListBooks :many
SELECT b.*, a.name as author_name
FROM books b
JOIN authors a ON b.author_id = a.id
ORDER BY b.title;

-- name: ListBooksByAuthor :many
SELECT * FROM books
WHERE author_id = ?
ORDER BY title;

-- name: UpdateBook :execresult
UPDATE books
SET title = ?,
    author_id = ?,
    published_year = ?,
    genre = ?,
    description = ?
WHERE id = ?;

-- name: DeleteBook :exec
DELETE FROM books
WHERE id = ?;