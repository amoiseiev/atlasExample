SET SEARCH_PATH TO shop;

CREATE SEQUENCE accounts_id_seq;
CREATE TABLE accounts
(
    id         INT       DEFAULT nextval('accounts_id_seq'),
    username   VARCHAR(50) UNIQUE                  NOT NULL,
    password   VARCHAR(50)                         NOT NULL,
    full_name  VARCHAR(100)                        NOT NULL,
    email      VARCHAR(255) UNIQUE                 NOT NULL,
    created_on TIMESTAMP DEFAULT current_timestamp NOT NULL,
    last_login TIMESTAMP                           NULL,
    PRIMARY KEY (id)
);

CREATE SEQUENCE states_id_seq;
CREATE TABLE states
(
    id   INT DEFAULT nextval('states_id_seq'),
    name VARCHAR(50) NOT NULL,
    PRIMARY KEY (id)
);

CREATE SEQUENCE addresses_id_seq;
CREATE TABLE addresses
(
    id             INT       DEFAULT nextval('addresses_id_seq'),
    account_id     INT                                 NOT NULL,
    address_1      VARCHAR(200)                        NOT NULL,
    address_2      VARCHAR(200)                        NULL,
    city           VARCHAR(50)                         NOT NULL,
    state_id       INT                                 NOT NULL,
    zip_code       VARCHAR(10)                         NOT NULL,
    recipient_name VARCHAR(100)                        NOT NULL,
    added_on       TIMESTAMP DEFAULT current_timestamp NOT NULL,
    PRIMARY KEY (id),
    FOREIGN KEY (account_id) REFERENCES accounts (id) ON DELETE CASCADE,
    FOREIGN KEY (state_id) REFERENCES states (id) ON DELETE CASCADE
);
CREATE INDEX addresses_account_id_idx ON addresses (account_id);
CREATE INDEX state_id_idx ON addresses (state_id);
