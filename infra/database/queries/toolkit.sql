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

-- name: InsertLoan :exec
INSERT INTO loans (id, toolkit_id, team_id, submitter_id, date, session_number, status)
VALUES ($1, $2, $3, $4, $5, $6, $7);

-- name: UpdateLoanNote :exec
UPDATE loans
SET
    note = $1,
    updated_at = $2
WHERE id = $3;

-- name: UpdateLoanStatus :exec
UPDATE loans
SET
    status = $1,
    updated_at = $2
WHERE id = $3;

-- name: FindLoansByTeamID :many
SELECT
    l.id, l.toolkit_id, t.kit_name,
    t.course_id, c.name AS course_name,
    l.team_id, l.submitter_id,
    l.date, l.session_number, l.status, l.note,
    l.created_at, l.updated_at
FROM loans l
LEFT JOIN toolkits t ON t.id = l.toolkit_id
LEFT JOIN repl_courses c ON c.id = t.course_id
WHERE l.team_id = $1
ORDER BY l.date, l.session_number;

-- name: FindLoansByDate :many
SELECT
    l.id, l.toolkit_id, t.kit_name,
    t.course_id, c.name AS course_name,
    l.team_id, l.submitter_id,
    l.date, l.session_number, l.status, l.note,
    l.created_at, l.updated_at
FROM loans l
LEFT JOIN toolkits t ON t.id = l.toolkit_id
LEFT JOIN repl_courses c ON c.id = t.course_id
WHERE l.date = $1
ORDER BY l.session_number;

-- name: FindLoanByID :one
SELECT
    l.id, l.toolkit_id, t.kit_name,
    t.course_id, c.name AS course_name,
    l.team_id, l.submitter_id,
    l.date, l.session_number, l.status, l.note,
    l.created_at, l.updated_at
FROM loans l
LEFT JOIN toolkits t ON t.id = l.toolkit_id
LEFT JOIN repl_courses c ON c.id = t.course_id
WHERE l.id = $1;
