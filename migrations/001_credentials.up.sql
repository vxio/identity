CREATE TABLE credentials (
    credential_id   VARCHAR(36) NOT NULL,
    identity_id     VARCHAR(36) NOT NULL,
    tenant_id       VARCHAR(36) NOT NULL,

    created_on      TIMESTAMP NOT NULL,
    last_used_on    TIMESTAMP NOT NULL,
    
    disabled_on     TIMESTAMP DEFAULT NULL,
    disabled_by     VARCHAR(36) DEFAULT NULL,

    CONSTRAINT credential_per_tenant UNIQUE (credential_id, tenant_id),
    CONSTRAINT credential_id_pk PRIMARY KEY (credential_id)
);
