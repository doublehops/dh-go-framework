CREATE TABLE my_new_table (
    id          CHAR(36) NOT NULL,
    currency_id CHAR(36) DEFAULT '',
    name        VARCHAR(100) DEFAULT '',
    created_by  CHAR(36) DEFAULT '',
    updated_by  CHAR(36) DEFAULT '',
    created_at  DATETIME,
    updated_at  DATETIME,
    deleted_at  DATETIME,
    PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
