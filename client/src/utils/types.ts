export type AuthResponse = {
  status: "ok" | "error";
  errorMessage: string;
};

export type AuthenticatorMetadata = {
  name: string;
  icon_dark: string;
  icon_light: string;
};

export type Passkey = {
  authenticator_metadata: AuthenticatorMetadata;
  credential_id: string;
  created_at: string;
  updated_at: string;
};
