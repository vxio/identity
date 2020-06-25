# Identity
**[Identity](README.md)** | **Configuration** | **[Running](RUNNING.md)** | **[Client](../pkg/client/README.md)**

---

## Configuration
Custom configuration for this application may be specified via an environment variable `APP_CONFIG` to a configuration file that will be merged with the default configuration file.

- [Default Configuration](../configs/config.default.yml)
- [Config Source Code](../pkg/identity/model_config.go)
- Full Configuration
  ```yaml
  Identity:

    # Service configurations
    Servers:

      # Public service configuration
      Public:
        Bind:
          # Address and port to listen on.
          Address: ":8200"

      # Health/Admin service configuration.
      Admin:
        Bind:
          # Address and port to listen on.
          Address: ":8201"

    # All database configuration is done here. Only one connector can be configured.
    Database:

      # Database name to use for selected connector.
      DatabaseName: "identity"

      # MySql configuration
      MySQL:  
        Address: tcp(mysqlidentity:3306)
        User: identity
        Password: identity

      # OR uses the sqllite db
      SQLLite:
        Path: ":memory:"

    # Handles all configuration for the invites functionality
    Invites:
    
      # Timeout the invites last for when created.
      Expiration: 48h

      # What host to place in the invite email to send the new user.
      SendToHost: https://api.moov.io 

      # What path to place in the invite email to send the new user.
      SendToPath: /authentication/tenants/{{.TenantID}}

    # Notifications configuration options for how to send email notifications like invites.
    # Only one of the sub configurations may be selected.
    Notifications:

      # Connect to an SMTP server.
      SMTP:
        Host: mailslurper
        Port: 1025
        User: test
        Pass: test
        From: noreply@moov.io
        SSL: true

      # Mock the sending of the email notification and just log it.
      Mock:
        from: noreply@moov.io

    # Gateway configuration to look up public keys to verify JWT tokens.
    Gateway:

      # If neither http or file are specified, the service will generate random keys
      Keys:

        # Pulls Keys from endpoints
        HTTP:
        URLs:
        - http://tumbler:8204/.well-known/jwks.json

        # Pulls keys from the disk
        File:
          Paths: 
          - ./configs/gateway-jwks-sig-pub.json

    # Authentication configuration for when the authentication service has finished verifying the user and sends them to identity
    Authentication:

      # Where to send the browser once they have been verified that they have an identity
      LandingUrl: "/whoami"

      # If neither http or file are specified, the service will generate random keys
      Keys:

        # Pulls Keys from endpoints
        HTTP:
          URLs:
          - http://authn:8202/.well-known/jwks.json

        # Pulls keys from the disk
        File:
          Paths: 
          - ./configs/authn-jwks-sig-pub.json

    # Session contains configuration about the JWT that is returned once they've `/authenticated` or `/registered` with identity. 
    # This is used as a longer term key to allow access to the API's that can be checked with the public key returned by hitting the 
    # <host>:<public port>/.well-known/jwks.json endpoint.
    Session:
      Expiration: 1h

      # If neither http or file are specified, the service will generate random keys. The public keys are available on the /.well-known/jwks.json endpoint
      Keys:
        Paths:
        - ./configs/identity-jwks-sig-pub.json
        - ./configs/identity-jwks-sig-priv.json

      # Note: HTTP isn't usable here due to needing the private key for signing. 
      # If you construct a separate service that doesn't use our webkeys package that outputs private keys it could be used here.
      # Public key verification in the webkeys package is verified on output to the endpoint not on ingestion.
      # HTTP:
      #   URLs:
      #   - http://customservice:5678/.well-known/jwks.json
  ```

---
**[Next - Running](RUNNING.md)**