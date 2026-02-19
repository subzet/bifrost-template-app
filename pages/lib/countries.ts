export const COUNTRIES: { code: string; name: string }[] = [
  { code: "AR", name: "Argentina" },
];

export function countryName(code: string): string {
  return COUNTRIES.find((c) => c.code === code)?.name ?? code;
}
