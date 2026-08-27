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

-- name: FindAllToolkits :many
SELECT * FROM toolkits
ORDER BY kit_name;

-- name: FindToolkitsWithInstruments :many
SELECT
    ti.toolkit_id,
    t.kit_name,
    t.course_id,
    t.total_count,
    t.out_of_order_count,
    t.created_at AS toolkit_created_at,
    t.updated_at AS toolkit_updated_at,
    ti.instrument_id,
    i.name AS instrument_name,
    i.created_at AS instrument_created_at,
    i.updated_at AS instrument_updated_at
FROM toolkit_instruments ti
INNER JOIN toolkits t ON ti.toolkit_id = t.id
INNER JOIN instruments i ON ti.instrument_id = i.id
ORDER BY t.kit_name, t.id;
