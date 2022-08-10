
/* name: SetConfig :exec */
INSERT INTO configs (uri, data)
VALUES (?1, ?2)
ON CONFLICT(uri) DO UPDATE
SET data = ?2;

/* name: GetConfig :one */
SELECT data FROM configs 
WHERE uri = ?;

/* name: Configs :many */
SELECT * FROM configs;

/* name: UnsetConfig :exec */
DELETE FROM configs WHERE uri=?;




/* name: CreateResource :one */
INSERT INTO resources (name, path, kind)
VALUES(?1, ?2, ?3)
RETURNING *;

/* name: FindResource :many */
SELECT * FROM resources;

/* name: FindResourceByKind :many */
SELECT * FROM resources
WHERE kind = ?;

/* name: FindResourceByNameByKind :one */
SELECT * FROM resources
WHERE name = ? AND kind = ?;

/* name: RemoveResource :exec */
DELETE FROM resources WHERE name=? AND kind = ?;

/* name: ActiveResource :one */
SELECT resources.name, resources.path FROM configs 
JOIN resources ON configs.data = resources.name 
WHERE configs.uri = ? AND resources.kind = ?;
