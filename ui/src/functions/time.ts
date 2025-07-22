
export const formatDateTime = (dateString: string): string => {
  if (!dateString) {
    return "";
  }
  // Truncate microseconds to milliseconds
  const truncatedDateString = dateString.replace(/(\.\d{3})\d+Z$/, '$1Z');
  const date = new Date(truncatedDateString);
  if (isNaN(date.getTime())) {
    return "Invalid Date";
  }

  const day = String(date.getDate()).padStart(2, '0');
  const month = String(date.getMonth() + 1).padStart(2, '0');
  const year = date.getFullYear();
  const hours = String(date.getHours()).padStart(2, '0');
  const minutes = String(date.getMinutes()).padStart(2, '0');
  const seconds = String(date.getSeconds()).padStart(2, '0');

  return `${day}/${month}/${year} ${hours}:${minutes}:${seconds}`;
};
