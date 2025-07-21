CREATE TABLE IF NOT EXISTS users
(
    id              SERIAL PRIMARY KEY,
    name            VARCHAR(255) NOT NULL,
    age             INTEGER      NOT NULL CHECK (age > 0),
    city            VARCHAR(255) NOT NULL,
    reg_dt          TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    gender          VARCHAR(1)   NOT NULL,
    search_gender   VARCHAR(1)   NOT NULL,
    search_age_from INTEGER      NOT NULL CHECK (search_age_from > 0),
    search_age_to   INTEGER      NOT NULL CHECK (search_age_to > search_age_from),
    location        FLOAT
);

-- Insert sample data
INSERT INTO users (name, age, city, gender, search_gender, search_age_from, search_age_to, location)
VALUES ('John Doe', 30, 'New York', 'm', 'f', 20, 30, 5),
       ('Jane Smith', 25, 'Los Angeles', 'f', 'm', 20, 30, 10),
       ('Bob Johnson', 35, 'Chicago', 'm', 'f', 18, 25, 15),
       ('Alice Williams', 28, 'Houston', 'f', 'm', 20, 30, 20),
       ('Charlie Brown', 32, 'Phoenix', 'f', 'f', 18, 40, 50),
       ('Diana Davis', 27, 'Philadelphia', 'm', 'f', 18, 80, 0),
       ('Eve Miller', 29, 'San Antonio', 'f', 'f', 18, 40, 100),
       ('Frank Wilson', 31, 'San Diego', 'm', 'f', 18, 40, 150),
       ('Grace Moore', 26, 'Dallas', 'f', 'f', 18, 40, 200),
       ('Henry Taylor', 33, 'San Jose', 'm', 'f', 18, 40, 250);