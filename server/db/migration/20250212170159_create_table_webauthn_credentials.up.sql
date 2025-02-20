CREATE TABLE webauthn_credentials (
  id UUID NOT NULL PRIMARY KEY,
  user_id UUID NOT NULL,
  credential_id BYTEA NOT NULL,
  public_key BYTEA NOT NULL,
  attestation_type VARCHAR(255) NOT NULL,
  transport VARCHAR(255),
  flags JSONB NOT NULL,
  authenticator JSONB NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
