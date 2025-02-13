SET
  statement_timeout = 0;
--bun:split
CREATE TABLE webauthn_credentials (
  id UUID NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id UUID NOT NULL,
  credential_id BYTEA NOT NULL,
  public_key BYTEA NOT NULL,
  attestation_type VARCHAR(255) NOT NULL,
  transport VARCHAR(255) [],
  flags JSONB NOT NULL,
  authenticator JSONB NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);
--bun:split
