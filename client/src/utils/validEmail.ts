export function isValidEmail(email: string): boolean {
  const emailRegex = /^[\w-\.]+@([\w-]+\.)+[\w-]{2,4}$/g; // eslint-disable-line
  return email !== "" && emailRegex.test(email);
}
