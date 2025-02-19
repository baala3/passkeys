export function isValidEmail(email: string): boolean {
  const emailRegex = /^[\w-\.]+@([\w-]+\.)+[\w-]{2,4}$/g;
  return email !== "" && emailRegex.test(email);
}
