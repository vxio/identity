CREATE TABLE identity (
    identity_id     VARCHAR(36) NOT NULL,
    tenant_id       VARCHAR(36) NOT NULL,

    first_name      VARCHAR(255) NOT NULL,
    middle_name     VARCHAR(255),
    last_name       VARCHAR(255) NOT NULL,
    nick_name       VARCHAR(255),
    suffix          VARCHAR(20),
    birth_date      TIMESTAMP,
    status          VARCHAR(20) NOT NULL,
    email           VARCHAR(255) NOT NULL,
    email_verified  BOOLEAN DEFAULT false,

    registered_on   TIMESTAMP NOT NULL,
    invite_id       VARCHAR(36),
    
    disabled_on     TIMESTAMP,
    disabled_by     VARCHAR(36),

    last_updated_on TIMESTAMP NOT NULL,
    
    CONSTRAINT identity_pk PRIMARY KEY (identity_id),
    CONSTRAINT email_tenant_uniq UNIQUE (tenant_id, email)
);