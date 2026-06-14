-- name: InsertToolkit :one
INSERT INTO toolkits (kit_name, total_count)
VALUES ($1, $2)
RETURNING id;

-- name: UpdateToolkitCourseID :exec
UPDATE toolkits
SET
    course_id = $1,
    updated_at = $2
WHERE id = $3;

-- name: InsertToolkitInstrument :exec
INSERT INTO toolkit_instruments (toolkit_id, instrument_id)
VALUES ($1, $2);

-- name: DeleteToolkitInstrument :exec
DELETE FROM toolkit_instruments
WHERE
        toolkit_id = $1
    AND
        instrument_id = $2;
