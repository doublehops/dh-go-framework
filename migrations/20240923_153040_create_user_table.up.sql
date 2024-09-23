CREATE TABLE organisation (
    id         INT(11) NOT NULL AUTO_INCREMENT,
    name       VARCHAR(255) NOT NULL,
    created_at DATETIME,
    updated_at DATETIME,
    deleted_at DATETIME,
    PRIMARY KEY (id)
);

------------------

CREATE TABLE user (
    id INT(11) NOT NULL AUTO_INCREMENT,
    organisation_id INT(11) NOT NULL,
    name VARCHAR(255) NOT NULL,
    email_address VARCHAR(255) NOT NULL,
    email_verified TINYINT(4),
    password VARCHAR(255) NOT NULL,
    password_reset_string VARCHAR(255) NOT NULL,
    password_reset_timeout DATETIME,
    is_active TINYINT(4),
    created_at DATETIME,
    updated_at DATETIME,
    deleted_at DATETIME,
    PRIMARY KEY (id),
    FOREIGN KEY (organisation_id) REFERENCES organisation(id)
);
