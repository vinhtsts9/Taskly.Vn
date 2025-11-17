-- +goose Up
-- +goose StatementBegin

-- =================================================================
-- Chèn dữ liệu cho Bảng `roles`
-- =================================================================
-- Sử dụng CTE (Common Table Expressions) để lưu trữ ID của các role vừa tạo
WITH new_roles (role_name) AS (
    VALUES
        ('buyer'),
        ('seller'),
        ('admin')
),
inserted_roles AS (
    INSERT INTO roles (role_name)
    SELECT role_name FROM new_roles
    RETURNING id, role_name
)
SELECT 'Roles inserted: ' || count(*) FROM inserted_roles;


-- =================================================================
-- Chèn dữ liệu cho Bảng `permissions`
-- =================================================================
-- Các quyền cơ bản cho các tài nguyên
WITH new_permissions (name, resource, action) AS (
    VALUES
        ('Đọc gigs', 'gigs', 'read'),
        ('Tạo gigs', 'gigs', 'create'),
        ('Cập nhật gigs', 'gigs', 'update'),
        ('Xóa gigs', 'gigs', 'delete'),
        ('Quản lý gigs', 'gigs', 'manage'),

        ('Đọc orders', 'orders', 'read'),
        ('Tạo orders', 'orders', 'create'),
        ('Cập nhật orders', 'orders', 'update'),

        ('Gửi tin nhắn', 'messages', 'send'),
        ('Đọc tin nhắn', 'messages', 'read'),

        ('Quản lý người dùng', 'users', 'manage'),
        ('Quản lý vai trò', 'roles', 'manage')
),
inserted_permissions AS (
    INSERT INTO permissions (name, resource, action)
    SELECT name, resource, action FROM new_permissions
    ON CONFLICT (resource, action) DO NOTHING
    RETURNING id, name
)
SELECT 'Permissions inserted: ' || count(*) FROM inserted_permissions;


-- =================================================================
-- Chèn dữ liệu cho Bảng `role_permissions` (Gán quyền cho vai trò)
-- =================================================================
-- Gán quyền cho vai trò 'admin'
INSERT INTO role_permissions (role_id, permission_id)
SELECT
    (SELECT id FROM roles WHERE role_name = 'admin'),
    p.id
FROM permissions p;

-- Gán quyền cho vai trò 'seller'
INSERT INTO role_permissions (role_id, permission_id)
SELECT
    (SELECT id FROM roles WHERE role_name = 'seller'),
    p.id
FROM permissions p
WHERE p.resource = 'gigs' OR (p.resource = 'messages' AND p.action = 'send') OR (p.resource = 'messages' AND p.action = 'read');

-- Gán quyền cho vai trò 'buyer'
INSERT INTO role_permissions (role_id, permission_id)
SELECT
    (SELECT id FROM roles WHERE role_name = 'buyer'),
    p.id
FROM permissions p
WHERE (p.resource = 'orders' AND p.action = 'create') OR (p.resource = 'orders' AND p.action = 'read') OR (p.resource = 'messages' AND p.action = 'send') OR (p.resource = 'messages' AND p.action = 'read');


-- =================================================================
-- Chèn dữ liệu người dùng mẫu và gán vai trò
-- Lưu ý: Bảng `users` không được định nghĩa trong context,
-- tôi giả định nó tồn tại và có cột `id UUID`.
-- Bạn cần tạo bảng `users` trước khi chạy migration này.
-- =================================================================
-- Ví dụ về bảng users:
-- CREATE TABLE users (
--     id UUID PRIMARY KEY DEFAULT uuid_generate_v1(),
--     username VARCHAR(50) UNIQUE NOT NULL,
--     email VARCHAR(255) UNIQUE NOT NULL,
--     password_hash VARCHAR(255) NOT NULL,
--     created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
--     updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
-- );

-- Tạo 3 người dùng mẫu
WITH sample_users (id, username, email, password_hash) AS (
    VALUES
        ('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11'::UUID, 'test_buyer', 'buyer@example.com', 'hashed_password_1'),
        ('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a12'::UUID, 'test_seller', 'seller@example.com', 'hashed_password_2'),
        ('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a13'::UUID, 'test_admin', 'admin@example.com', 'hashed_password_3')
)
-- Giả sử bạn có bảng users, bạn sẽ INSERT vào đó.
-- INSERT INTO users (id, username, email, password_hash)
-- SELECT id, username, email, password_hash FROM sample_users;

-- Gán vai trò cho người dùng mẫu
INSERT INTO user_roles (user_id, role_id)
VALUES
    ('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', (SELECT id FROM roles WHERE role_name = 'buyer')),
    ('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a12', (SELECT id FROM roles WHERE role_name = 'seller')),
    ('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a13', (SELECT id FROM roles WHERE role_name = 'admin'));


-- =================================================================
-- Chèn dữ liệu cho Bảng `rooms` và `messages`
-- =================================================================
-- Tạo một phòng chat giữa buyer và seller
WITH new_room AS (
    INSERT INTO rooms (user1_id, user2_id)
    VALUES ('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a12')
    RETURNING id
)
-- Chèn một vài tin nhắn vào phòng chat vừa tạo
INSERT INTO messages (room_id, sender_id, receiver_id, content)
VALUES
    ((SELECT id FROM new_room), 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a12', 'Chào bạn, tôi quan tâm đến dịch vụ của bạn.'),
    ((SELECT id FROM new_room), 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a12', 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'Chào bạn. Rất vui được hỗ trợ. Bạn cần gì ạ?');


-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- Xóa dữ liệu theo thứ tự ngược lại để đảm bảo không vi phạm khóa ngoại
DELETE FROM messages WHERE sender_id IN ('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a12');
DELETE FROM rooms WHERE user1_id = 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11';
DELETE FROM user_roles WHERE user_id IN ('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a12', 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a13');
-- DELETE FROM users WHERE username IN ('test_buyer', 'test_seller', 'test_admin'); -- Bỏ comment nếu bạn có bảng users
DELETE FROM role_permissions;
DELETE FROM permissions;
DELETE FROM roles;

-- +goose StatementEnd
