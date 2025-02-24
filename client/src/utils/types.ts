export type AuthResponse = {
  status: "ok" | "error";
  errorMessage: string;
};

export type PasskeyProvider = {
  name: string;
  icon_dark: string;
  icon_light: string;
};

// for passkey management page
export type Passkey = {
  passkey_provider: PasskeyProvider;
  credential_id: string;
  created_at: string;
  updated_at: string;
};
