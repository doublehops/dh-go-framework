CREATE TABLE organisation (
    id         CHAR(36) NOT NULL,
    name       VARCHAR(255) NOT NULL DEFAULT '',
    created_by CHAR(36) DEFAULT '',
    updated_by CHAR(36) DEFAULT '',
    created_at DATETIME,
    updated_at DATETIME,
    deleted_at DATETIME,
    PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Matches userservice.DefaultOrganisationID.
INSERT INTO organisation (id, name) VALUES ('00000000-0000-0000-0000-000000000100', 'org1');

-- Global (system-level) roles. IDs match the constants in internal/model/userrole.
CREATE TABLE user_role (
    id         CHAR(36) NOT NULL,
    name       VARCHAR(100) NOT NULL DEFAULT '',
    created_by CHAR(36) DEFAULT '',
    updated_by CHAR(36) DEFAULT '',
    created_at DATETIME,
    updated_at DATETIME,
    deleted_at DATETIME,
    PRIMARY KEY (id),
    UNIQUE KEY uq_user_role_name (name)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

INSERT INTO user_role (id, name) VALUES
    ('6649c54a-4641-4554-9ec3-497192ac3df3', 'standard_user'),
    ('6a8c21d8-201d-4420-a964-54a2b99c3db9', 'organisation_admin'),
    ('4ad2b30f-81d7-4b58-b950-db5c51e9590e', 'master_admin');

CREATE TABLE user (
    id                     CHAR(36) NOT NULL,
    organisation_id        CHAR(36) NOT NULL,
    role_id                CHAR(36) NOT NULL DEFAULT '6649c54a-4641-4554-9ec3-497192ac3df3',
    name                   VARCHAR(255) NOT NULL DEFAULT '',
    email_address          VARCHAR(255) NOT NULL DEFAULT '',
    email_verified         TINYINT(4) DEFAULT 0,
    email_verified_token   VARCHAR(255) NULL DEFAULT '',
    password               VARCHAR(255) NOT NULL DEFAULT '',
    password_reset_token   VARCHAR(255) NULL DEFAULT '',
    password_reset_expire  DATETIME,
    is_active              TINYINT(4) DEFAULT 0,
    created_by             CHAR(36) DEFAULT '',
    updated_by             CHAR(36) DEFAULT '',
    created_at             DATETIME,
    updated_at             DATETIME,
    deleted_at             DATETIME,
    PRIMARY KEY (id),
    FOREIGN KEY (organisation_id) REFERENCES organisation(id),
    FOREIGN KEY (role_id) REFERENCES user_role(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
