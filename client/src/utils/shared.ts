function isValidEmail(email: string): boolean {
  const emailRegex = /^[\w-\.]+@([\w-]+\.)+[\w-]{2,4}$/g; // eslint-disable-line
  return email !== "" && emailRegex.test(email);
}

function isAuthenticated(): boolean {
  return document.cookie.indexOf("auth") !== -1;
}

function logout() {
  document.cookie = "auth=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/;";
}

export { isValidEmail, isAuthenticated, logout };
