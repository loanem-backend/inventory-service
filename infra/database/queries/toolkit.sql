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
