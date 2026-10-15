CREATE TABLE author
(
    id         CHAR(36) NOT NULL,
    name       VARCHAR(100) NOT NULL DEFAULT '',
    created_by CHAR(36) DEFAULT '',
    updated_by CHAR(36) DEFAULT '',
    created_at DATETIME,
    updated_at DATETIME,
    deleted_at DATETIME,
    PRIMARY KEY (id)
) ENGINE = InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
