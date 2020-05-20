CREATE TABLE invites (
    invite_id       VARCHAR(36) NOT NULL,
    tenant_id       VARCHAR(36) NOT NULL,

    email           VARCHAR(255) NOT NULL,
    invited_by      VARCHAR(36) NOT NULL,
    invited_on      TIMESTAMP NOT NULL,
    redeemed_on     TIMESTAMP,
    expires_on      TIMESTAMP NOT NULL,
    disabled_on     TIMESTAMP DEFAULT NULL,
    disabled_by     VARCHAR(36) DEFAULT NULL,

    secret_code     VARCHAR(255) NOT NULL,

    CONSTRAINT invite_id_pk PRIMARY KEY (invite_id)
);