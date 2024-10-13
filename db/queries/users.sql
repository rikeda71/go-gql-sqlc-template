-- name: InsertUser :one
INSERT INTO users /* users_001 */
(id, user_name, email) VALUES ($1, $2, $3)
RETURNING *;

-- name: FindUserByID :one
SELECT /* users_002 */
    *
FROM users
WHERE id = $1;
