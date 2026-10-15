CREATE TABLE user_session (
    id           CHAR(36) NOT NULL,
    user_id      CHAR(36) NOT NULL,
    token        CHAR(44) DEFAULT '',
    last_request DATETIME,
    created_by   CHAR(36) DEFAULT '',
    updated_by   CHAR(36) DEFAULT '',
    created_at   DATETIME,
    updated_at   DATETIME,
    deleted_at   DATETIME,
    PRIMARY KEY (id),
    FOREIGN KEY (user_id) REFERENCES user(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
