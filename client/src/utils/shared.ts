import DateObject from "react-date-object";

function isValidEmail(email: string): boolean {
  const emailRegex = /^[\w-\.]+@([\w-]+\.)+[\w-]{2,4}$/g; // eslint-disable-line
  return email !== "" && emailRegex.test(email);
}

function formatDate(date: string) {
  const dateObject = new DateObject({ date: new Date(date).toISOString() });
  return dateObject.format("DD MMM YY, hh:mm a");
}

export { isValidEmail, formatDate };
