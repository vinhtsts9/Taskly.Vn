-- +goose Up
-- +goose StatementBegin

-- Bảng `roles` để lưu các vai trò (ví dụ: admin, buyer, seller)
CREATE TABLE roles (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v1(),
    role_name VARCHAR(10) CHECK (role_name IN ('buyer', 'seller', 'admin')) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Bảng `permissions` để lưu các quyền hạn chi tiết
-- Ví dụ: có quyền 'tạo' (action) trên tài nguyên 'gigs' (resource)
CREATE TABLE permissions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v1(),
    name VARCHAR(255) NOT NULL, -- Tên quyền để dễ nhận biết, ví dụ: "Tạo gig mới"
    resource VARCHAR(100) NOT NULL, -- Tên tài nguyên, ví dụ: "gigs", "orders", "users"
    action VARCHAR(50) NOT NULL, -- Hành động, ví dụ: "create", "read", "update", "delete", "manage"
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(resource, action)
);

-- Bảng `user_roles` là bảng trung gian để gán vai trò cho người dùng
-- Mối quan hệ nhiều-nhiều giữa `users` và `roles`
CREATE TABLE user_roles (
    user_id UUID NOT NULL,
    role_id UUID NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (user_id, role_id)
);

-- Bảng `role_permissions` là bảng trung gian để gán quyền cho vai trò
-- Mối quan hệ nhiều-nhiều giữa `roles` và `permissions`
CREATE TABLE role_permissions (
    role_id UUID NULL,
    permission_id UUID NOT NULL,
    PRIMARY KEY (role_id, permission_id)
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS role_permissions;
DROP TABLE IF EXISTS user_roles;
DROP TABLE IF EXISTS permissions;
DROP TABLE IF EXISTS roles;

-- +goose StatementEnd 