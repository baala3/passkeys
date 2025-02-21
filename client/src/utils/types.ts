export type AuthResponse = {
  status: "ok" | "error";
  errorMessage: string;
};

export type Passkey = {
  aaguid: string;
  credential_id: string;
  created_at: string;
  updated_at: string;
};
