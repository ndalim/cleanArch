DROP TABLE IF EXISTS users;
CREATE TABLE users (
    id INT NOT NULL AUTO_INCREMENT,
    name VARCHAR(100),
    hashpass VARCHAR(100),
    PRIMARY KEY (id)
);

INSERT INTO users (name) VALUES
('alex'),
('val');