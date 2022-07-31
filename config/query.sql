/* name: AddToolLibrary :one */
INSERT INTO resources (name, path, kind)
VALUES(?1, ?2, "tool.library")
RETURNING *;

/* name: ActiveToolLibrary :one */
SELECT resources.name, resources.path FROM configs 
JOIN resources ON configs.data = resources.name 
WHERE configs.uri = "library.active";

/* name: SetActiveLibrary :exec */
INSERT INTO configs (uri, data) 
VALUES("library.active", ?1)
ON CONFLICT(uri) DO UPDATE 
SET data = ?1;

