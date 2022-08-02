/* name: AddToolLibrary :one */
INSERT INTO resources (name, path, kind)
VALUES(?1, ?2, "tool.library")
RETURNING *;

/* name: FindResourceByName :one */
SELECT * FROM resources
WHERE name = ?;

/* name: ToolLibraries :many */
SELECT * FROM resources
WHERE kind = "tool.library";

/* name: RemoveResource :exec */
DELETE FROM resources WHERE name=?;

/* name: ActiveToolLibrary :one */
SELECT resources.name, resources.path FROM configs 
JOIN resources ON configs.data = resources.name 
WHERE configs.uri = "library.active";

/* name: SetActiveLibrary :exec */
INSERT INTO configs (uri, data) 
VALUES("library.active", ?1)
ON CONFLICT(uri) DO UPDATE 
SET data = ?1;

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
