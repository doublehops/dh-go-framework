CREATE TABLE organisation (
    id         INT(11) NOT NULL AUTO_INCREMENT,
    name       VARCHAR(255) NOT NULL,
    created_at DATETIME,
    updated_at DATETIME,
    deleted_at DATETIME,
    PRIMARY KEY (id)
);

------------------

INSERT INTO organisation (name) VALUES ('org1')
------------------

CREATE TABLE user (
    id INT(11) NOT NULL AUTO_INCREMENT,
    organisation_id INT(11) NOT NULL,
    name VARCHAR(255) NOT NULL,
    email_address VARCHAR(255) NOT NULL,
    email_verified TINYINT(4) DEFAULT 0,
    email_verified_token VARCHAR(255) NULL DEFAULT "",
    password VARCHAR(255) NOT NULL,
    password_reset_token VARCHAR(255) NULL DEFAULT "",
    password_reset_expire DATETIME,
    is_active TINYINT(4),
    created_at DATETIME,
    updated_at DATETIME,
    deleted_at DATETIME,
    PRIMARY KEY (id),
    FOREIGN KEY (organisation_id) REFERENCES organisation(id)
);
