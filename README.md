An implementation of authentication with passkeys using [go-webauthn](https://github.com/go-webauthn/webauthn)
library.

## Local setup

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

This is to test password managers like bitwarden.
Create a custom HTTPS URL that will route traffic to your your local server. In this case `http://localhost:9044`

```powershell
 ngrok http http://localhost:9044
```

Look for the Forwarding output.

```powershell
Forwarding      https://51ed-47-150-126-75.ngrok-free.app -> http://localhost:9044
```

Fix up the .env file.

```env
RP_DISPLAY_NAME=PasskeyDemo
RP_ID=51ed-47-150-126-75.ngrok-free.app
RP_ORIGIN=https://51ed-47-150-126-75.ngrok-free.app
```
