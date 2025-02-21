export type AuthResponse = {
  status: "ok" | "error";
  errorMessage: string;
};

export type Passkey = {
  aaguid: string;
  sign_count: number;
  created_at: string;
  updated_at: string;
};
