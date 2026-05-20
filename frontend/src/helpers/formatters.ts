const getString = (value: unknown, fallback: string = ""): string =>
  typeof value === "string" ? value : fallback;

const getNumber = (value: unknown, fallback: number = 0): number => {
  const num = Number(value);
  return isNaN(num) ? fallback : num;
};

const formatDate = (value: unknown): string => {
  if (typeof value !== "string") return "";
  const date = new Date(value);
  return date.toISOString().split("T")[0] ?? "";
};

const getDateToString = (value: unknown): string => {
  const str = getString(value);
  if (!str) return "";

  const date = new Date(str);
  if (isNaN(date.getTime())) return "";

  return date.toISOString().split("T")[0] ?? "";
};

export { getString, getNumber, formatDate, getDateToString };
