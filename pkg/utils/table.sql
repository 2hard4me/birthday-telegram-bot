CREATE TABLE birthdays(
    id SERIAL NOT NULL,
    chat_id BIGINT NOT NULL,
    name TEXT NOT NULL,
    day INT NOT NULL,
    month INT NOT NULL,
    PRIMARY KEY (id)
)
