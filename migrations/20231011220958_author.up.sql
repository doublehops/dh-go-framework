CREATE TABLE author
(
    id         INT(11) AUTO_INCREMENT NOT NULL,
    user_id    INT(11) NOT NULL,
    name       varchar(100) NOT NULL,
    created_by INT(11),
    updated_by INT(11),
    created_at DATETIME,
    updated_at DATETIME,
    deleted_at DATETIME,
    PRIMARY KEY (id)
) ENGINE = InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
	