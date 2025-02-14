An implementation of authentication with passkeys using [go-webauthn](https://github.com/go-webauthn/webauthn)
library.

## Local setup

1. Start Postgres and Redis with Docker Compose:

```powershell
docker compose up -d
```

2. Setup the database:

```powershell
cd server
go run ./db/migration db init
go run ./db/migration db migrate
```

3. Start the server:

```powershell
go run .
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
