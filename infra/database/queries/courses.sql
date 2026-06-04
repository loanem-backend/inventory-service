-- name: InsertCourse :exec
INSERT INTO repl_courses (id, name, year)
VALUES ($1, $2, $3);

-- name: UpdateCourseName :exec
UPDATE repl_courses
SET
    name = $1,
    updated_at = $2
WHERE id = $3;

-- name: DeleteCourseByID :exec
DELETE FROM repl_courses
WHERE id = $1;
