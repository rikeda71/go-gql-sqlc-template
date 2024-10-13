-- migrate:up
CREATE TABLE users (
    id CHAR(36) PRIMARY KEY,
    user_name VARCHAR(50) NOT NULL UNIQUE,
    email VARCHAR(100) NOT NULL UNIQUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
COMMENT ON TABLE users IS 'ユーザ情報';
COMMENT ON COLUMN users.id IS 'ユーザID';
COMMENT ON COLUMN users.user_name IS 'ユーザ名';
COMMENT ON COLUMN users.email IS 'メールアドレス';
COMMENT ON COLUMN users.created_at IS '作成日';
COMMENT ON COLUMN users.updated_at IS '更新日';

-- migrate:down
DROP TABLE users;
