-- name: RecordHistory :exec
INSERT INTO history (command, cwd, timestamp, exit_code, duration_ms)
VALUES (?, ?, ?, ?, ?);

-- name: GetHistory :many
SELECT * FROM history
ORDER BY timestamp DESC
LIMIT ?;

-- name: SearchHistory :many
SELECT * FROM history
WHERE command LIKE ?
ORDER BY timestamp DESC
LIMIT ?;
