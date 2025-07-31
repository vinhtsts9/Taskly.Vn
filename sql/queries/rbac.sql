-- name: CreateRole :one
INSERT INTO roles (name)
VALUES ($1)
RETURNING *;

-- name: GetRoleByName :one
SELECT * FROM roles
WHERE name = $1;

-- name: AddRoleToUser :exec
INSERT INTO user_roles (user_id, role_id)
VALUES ($1, $2);

-- name: GetRolesByUserID :many
SELECT r.*
FROM roles r
JOIN user_roles ur ON r.id = ur.role_id
WHERE ur.user_id = $1;

-- name: RemoveRoleFromUser :exec
DELETE FROM user_roles
WHERE user_id = $1 AND role_id = $2;

-- name: CreatePermission :one
INSERT INTO permissions (name, resource, action)
VALUES ($1, $2, $3)
RETURNING *;

-- name: AddPermissionToRole :exec
INSERT INTO role_permissions (role_id, permission_id)
VALUES ($1, $2);

-- name: RemovePermissionFromRole :exec
DELETE FROM role_permissions
WHERE role_id = $1 AND permission_id = $2;

-- name: GetPermissionsByRoleID :many
SELECT p.*
FROM permissions p
JOIN role_permissions rp ON p.id = rp.permission_id
WHERE rp.role_id = $1;

-- name: GetPermissionsByUserID :many
SELECT p.*
FROM permissions p
JOIN role_permissions rp ON p.id = rp.permission_id
JOIN user_roles ur ON rp.role_id = ur.role_id
WHERE ur.user_id = $1; 