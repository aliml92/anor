-- name: GetUserByEmail :one
SELECT * FROM users WHERE email=$1;

-- name: CreateUser :exec
INSERT INTO users 
    (email, password, full_name)
VALUES 
    ($1, $2, $3);

-- name: UpdateUserStatusByEmail :exec
UPDATE users
SET status = $1
WHERE email = $2;

-- name: GetUser :one
SELECT * FROM users WHERE id=$1; 

-- name: CreateSeller :one
INSERT INTO users (
    email, password, full_name, status
) VALUES (
    $1, $2, $3, $4
) RETURNING id;


-- name: UpdateUserPassword :exec
UPDATE users
SET password = $1
WHERE id = $2 AND status = 'Active';