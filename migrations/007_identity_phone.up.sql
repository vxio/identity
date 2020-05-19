CREATE TABLE identity_phone (
    identity_id     VARCHAR(36) NOT NULL,
    phone_id        VARCHAR(36) NOT NULL,

    type            VARCHAR(20) NOT NULL,
    number          VARCHAR(15),
    validated       BOOLEAN NOT NULL,

    last_updated_on TIMESTAMP NOT NULL,

    CONSTRAINT phone_pk PRIMARY KEY (phone_id)
);