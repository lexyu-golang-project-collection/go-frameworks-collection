-- name: GetBook :one
SELECT * FROM book
WHERE id = ? LIMIT 1;

-- name: ListBooks :many
SELECT * FROM book
ORDER BY title;

-- name: CreateBook :one
INSERT INTO book (
  title, author
) VALUES (
  ?, ?
)
RETURNING *;

-- name: UpdateBook :exec
UPDATE book
set title = ?,
    author = ?
WHERE id = ?;

-- name: DeleteBook :exec
DELETE FROM book
WHERE id = ?;

-- name: SearchBooksByTitle :many
SELECT * FROM book
WHERE title LIKE '%' || ? || '%' COLLATE NOCASE;

-- name: SearchBooksByAuthor :many
SELECT * FROM book
WHERE author LIKE '%' || ? || '%' COLLATE NOCASE;

-- name: GetBooksPaginated :many
SELECT * FROM book 
ORDER BY id
LIMIT ? OFFSET ?;

-- name: CountBooks :one
SELECT COUNT(*) as total FROM book;

-- name: GetBooksByAuthor :many
SELECT * FROM book
WHERE author = ? COLLATE NOCASE
ORDER BY title;