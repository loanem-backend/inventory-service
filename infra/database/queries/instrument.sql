-- name: InsertInstrument :one
INSERT INTO instruments (name)
VALUES ($1)
RETURNING id;
