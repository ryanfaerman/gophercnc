/* name: GetAuthor :one */ 
SELECT * FROM authors
WHERE id = ? LIMIT 1;

/* name: listauthors :many */
SELECT * FROM authors
ORDER BY name;

/* name: CreateAuthor :one */
INSERT INTO authors (
  name, bio
) VALUES (
  ?1, ?2
) 
RETURNING id;

/* name: DeleteAuthor :exec */
DELETE FROM authors
WHERE id = ?;

/* name: ActiveLibraries :many */
SELECT * FROM libraries 
WHERE active = true;
