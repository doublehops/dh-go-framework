CREATE TABLE my_new_table (
    id INT(11) NOT NULL AUTO_INCREMENT,
    currency_id INT(11),
    name VARCHAR(100),
    created_at DATETIME,
    updated_at DATETIME,
    deleted_at DATETIME,
    PRIMARY KEY (id)
);
