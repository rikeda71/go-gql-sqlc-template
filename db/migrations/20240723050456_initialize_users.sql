-- migrate:up
CREATE TABLE users (
    id CHAR(36) PRIMARY KEY,
    user_name VARCHAR(50) NOT NULL UNIQUE,
    email VARCHAR(100) NOT NULL UNIQUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
COMMENT ON TABLE users IS 'User Information';
COMMENT ON COLUMN users.id IS 'UserID';
COMMENT ON COLUMN users.user_name IS 'UserName';
COMMENT ON COLUMN users.email IS 'Email Address';
COMMENT ON COLUMN users.created_at IS 'Creation Date';
COMMENT ON COLUMN users.updated_at IS 'Last Update Date';

-- migrate:down
DROP TABLE users;
