CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    surname VARCHAR(255),
    female VARCHAR(10),
    email VARCHAR(255),
    status_id INTEGER REFERENCES user_status(id),
    rule_id INTEGER REFERENCES rules(id)
);
