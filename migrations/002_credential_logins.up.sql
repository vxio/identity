CREATE TABLE credential_logins (
    credential_id   VARCHAR(36) NOT NULL,
    nonce           VARCHAR(255) NOT NULL,
    ip              VARCHAR(15) NOT NULL,
    created_on      TIMESTAMP NOT NULL,

    CONSTRAINT login_pk PRIMARY KEY (credential_id, nonce)
);