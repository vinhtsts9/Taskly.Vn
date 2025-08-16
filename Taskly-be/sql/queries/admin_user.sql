-- name: AdminListUsers :many
SELECT 
    u.id,
    u.names,
    COALESCE(u.profile_pic,'') as profile_pic,
    u.created_at,
    b.email,
    b.states,
    COALESCE(r.role_name,'')
FROM users u
JOIN user_base b ON u.user_base_id = b.user_base_id 
LEFT JOIN user_roles ur ON ur.user_id = u.id
LEFT JOIN roles r ON r.id = ur.role_id
WHERE ($1 = '' OR LOWER(u.names) LIKE LOWER('%' || $1 || '%') OR LOWER(b.email) LIKE LOWER('%' || $1 || '%'))
ORDER BY u.created_at DESC
LIMIT $2 OFFSET $3;

-- name: AdminCountUsers :one
SELECT COUNT(*)
FROM users u
JOIN user_base b ON u.user_base_id = b.user_base_id
WHERE ($1 = '' OR LOWER(u.names) LIKE LOWER('%' || $1 || '%') OR LOWER(b.email) LIKE LOWER('%' || $1 || '%'));

