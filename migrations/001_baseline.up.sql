CREATE TABLE credentials (
    credential_id   UUID NOT NULL,
    provider        UUID NOT NULL,
    subject_id      UUID NOT NULL,
    identity_id     UUID NOT NULL,
    created_on      TIMESTAMP NOT NULL,
    last_used_on    TIMESTAMP NOT NULL,
    disabled_on     TIMESTAMP DEFAULT NULL,
    disabled_by     UUID DEFAULT NULL,

    CONSTRAINT unique_login UNIQUE (provider, subject_id),
    CONSTRAINT credential_id_pk PRIMARY KEY (credential_id)
);

CREATE TABLE invites (
    invite_id       UUID NOT NULL,
    tenant_id       UUID NOT NULL,

    email           VARCHAR(255) NOT NULL,
    invited_by      UUID NOT NULL,
    invited_on      TIMESTAMP NOT NULL,
    redeemed_on     TIMESTAMP,
    expires_on      TIMESTAMP NOT NULL,
    disabled_on     TIMESTAMP DEFAULT NULL,
    disabled_by     UUID DEFAULT NULL,

    secret_code     VARCHAR(255) NOT NULL,

    CONSTRAINT invite_id_pk PRIMARY KEY (invite_id)
);

CREATE TABLE identity (
    identity_id     UUID NOT NULL,
    tenant_id       UUID NOT NULL,

    first_name      VARCHAR(255) NOT NULL,
    middle_name     VARCHAR(255),
    last_name       VARCHAR(255) NOT NULL,
    nick_name       VARCHAR(255),
    suffix          VARCHAR(20),
    birth_date      DATE NOT NULL,
    status          VARCHAR(20) NOT NULL,
    email           VARCHAR(255) NOT NULL,
    email_verified  BOOLEAN DEFAULT false,

    registered_on   TIMESTAMP NOT NULL,
    invite_id       UUID NOT NULL,
    
    disabled_on     TIMESTAMP,
    disabled_by     UUID,

    last_updated_on TIMESTAMP NOT NULL,
    
    CONSTRAINT identity_pk PRIMARY KEY (identity_id),
    CONSTRAINT email_tenant_uniq UNIQUE (tenant_id, email)
);

CREATE TABLE identity_address (
    identity_id     UUID NOT NULL,
    address_id      UUID NOT NULL,

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

CREATE INDEX identity_address_identity_id ON identity_address (identity_id);

CREATE TABLE identity_phone (
    identity_id     UUID NOT NULL,
    phone_id        UUID NOT NULL,

    type            VARCHAR(20) NOT NULL,
    number          VARCHAR(15),
    validated       BOOLEAN NOT NULL,

    last_updated_on TIMESTAMP NOT NULL,

    CONSTRAINT phone_pk PRIMARY KEY (phone_id)
);

CREATE INDEX identity_phone_identity_id ON identity_phone (identity_id);