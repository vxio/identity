# Identity
**Purpose** | **[Configuration](CONFIGURATION.md)** | **[Running](RUNNING.md)** | **[Client](../pkg/client/README.md)**

---

## Purpose

Identity is one of a trio of services that handle authentication in the system. Identity's purpose is to manage the details of the user and finalize the authentication flow once they have verified the person logging in is a member of the system. 

Once the identity has been authenticated or registered the user will be given a JWT that lasts for a period of time. This JWT is only usable on the domain they logged into and only from their same IP.

Users are authenticated by a service like [`authn`](https://github.com/moov-io/authn) by providing a token to the browser that they authenticated with against a provider. The authenticating service forwards the browser to either `/authenticated` or `/register` where the token is read and then the credentials within it are checked against what they registered their user.

## Invites

New identities can only be invited to join by another identity. The sender of an invite will call the invite endpoint with an email address of the person to invite. The receiptient of the invite will have to click a link from the email which contains an invite code that will forward them to a login screen specific to the Tenant of the Inviter. Once they authenticate with the method chosen by the tenant they will be forwarded to Identity to register their new Identity with as much pre-filled for them that we could obtain from the authentication. Identity requires a lot of user information including addresses and phone numbers so that we can run KYC/OFAC checks on them. Once they register they will be logged into the system.

## Credentials

Credentials are the ways a identity is allowed to log into the system. A credential is composed of a few values that must match up to the authentication method values. `Provider` and `SubjectID` work as a pair of values that uniquely identity. `Provider` is the method to which they authenticated with. `SubjectID` is the unique ID of that identity that the `Provider` must provide. These two values act as like a username and password to log into the system. These are safe as the system that generates the token (say AuthN) must ensure that the token cannot be tampered with and Identity can verify that the token was constructed by this other system. For example in Identity this verification is handled by a RSA public and private key combo. AuthN is the only one who has the private key to signed the token. Identity is able to obtain the public key from the AuthN service to be able to verify the token. We also go a step further to protect PII data for a registering user and encrypt the token.

Credential may be disabled and a disabled credential is no longer usable.

## Identities

Identities are never completed removed from the system. This is due to being able to reconstruct the history of this user for audit purposes. Identities can be disabled which will not allow them access into the system any longer and be hidden.

The trail of how the identity came to be in the system can be followed through the `InviteID` attribute of the identity. Every identity came in via a invite to join and every invite contains the `IdentityID` of the person who sent it and when.

---
**[Next - Configuration](CONFIGURATION.md)**
