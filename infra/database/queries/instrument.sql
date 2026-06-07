-- name: InsertInstrument :one
INSERT INTO instruments (name)
VALUES ($1)
RETURNING id;

-- name: DeleteInstrumentByID :exec
DELETE FROM instruments
WHERE id = $1;

-- name: UpdateInstrumentName :exec
UPDATE instruments
SET
    name = $1,
    updated_at = $2
WHERE id = $3;

-- name: UpdateInstrumentPicture :exec
UPDATE instruments
SET
    picture = $1,
    updated_at = $2
WHERE id = $3;

-- name: FindAllInstruments :many
SELECT * FROM instruments
ORDER BY name;

-- name: FindInstrumentsByToolkitID :many
SELECT i.* FROM instruments i
INNER JOIN toolkit_instruments ti ON ti.instrument_id = i.id
WHERE ti.toolkit_id = $1
ORDER BY name;
