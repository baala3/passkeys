An implementation of passwordless authentication using WebAuthn, including multiple authentication methods and security controls.

## Overview

This project shows how to use passkeys for smooth and secure authentication.
You can compare registeration and authentication flows using password and passkey.

- [Read my blog post on Passkey Introduction and Auth Flows](https://balashekhar-blog.netlify.app/blogs/passkey-introduction/passkey-introduction/)
- Passkey flows you can check here https://youtu.be/wFbZHGvNUzo

## Features

#### 1. Authentication Methods
  - Passkey Registration & Login - authentication using WebAuthn
  - Passkey Autofill - Browser supported credential autofill
  - Traditional Password Auth - Email/password as a fallback option
  - Discoverable Login - Sign in without typing a username

#### 2. Security Features
  - 2FA for Sensitive Actions, Passkey auth required for
    - Password changes
    - Email changes
    - Account deletion
  - Session Management (using Redis)
  - Credential Management (To add/remove passkeys)
  - AAGUID (to get authenticator data, ref [here](https://github.com/passkeydeveloper/passkey-authenticator-aaguids))

#### 3. Technical Stack
  - Frontend: React 19, TypeScript, Vite, Tailwind CSS, [SimpleWebAuthn](https://www.npmjs.com/package/@simplewebauthn/browser)
  - Backend: Go, Echo framework, [go-webauthn](https://github.com/go-webauthn/webauthn)
  - Database: PostgreSQL with Bun ORM
  - Cache: Redis for session storage
  - Deployment: Docker Compose with multi-stage builds

---

## Local setup
note: needs docker and docker compose.

1. setup the app (install dependencies and initialize the database schema)

```powershell
make install
```

2. start the app (starts the Go server together with Redis and Postgres)

```powershell
make run
```

---

**HTTPS setup**

To test with password managers like Bitwarden, set up HTTPS tunneling:

1. Create a custom HTTPS URL that will route traffic to your local server (i.e `http://localhost:9044`)

```shell
 ngrok http http://localhost:9044
```

2. Look for the Forwarding URL output.

```shell
Forwarding      https://51ed-47-150-126-75.ngrok-free.app -> http://localhost:9044
```

3. update .env file.

```shell
RP_DISPLAY_NAME=PasskeyDemo
RP_ID=51ed-47-150-126-75.ngrok-free.app
RP_ORIGIN=https://51ed-47-150-126-75.ngrok-free.app
```
