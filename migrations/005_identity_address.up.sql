CREATE TABLE identity_address (
    identity_id     VARCHAR(36) NOT NULL,
    address_id      VARCHAR(36) NOT NULL,

    type            VARCHAR(20) NOT NULL,
    address_1       VARCHAR(255) NOT NULL,
    address_2       VARCHAR(255),
    city            VARCHAR(255) NOT NULL,
    state           VARCHAR(2) NOT NULL,
    country         VARCHAR(2) NOT NULL,
    validated       BOOLEAN NOT NULL,

    last_updated_on TIMESTAMP NOT NULL,

    CONSTRAINT address_pk PRIMARY KEY (address_id)
);