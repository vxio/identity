CREATE TABLE identity (
    identity_id     VARCHAR(36) NOT NULL,

    first_name      VARCHAR(255) NOT NULL,
    middle_name     VARCHAR(255),
    last_name       VARCHAR(255) NOT NULL,
    nick_name       VARCHAR(255),
    suffix          VARCHAR(20),
    birth_date      TIMESTAMP,
    email           VARCHAR(255) NOT NULL,
    email_verified  BOOLEAN DEFAULT false,
    last_updated_on TIMESTAMP NOT NULL,
    
    CONSTRAINT identity_pk PRIMARY KEY (identity_id)
);