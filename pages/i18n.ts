export function t(
  translations: Record<string, string>,
  key: string,
  params?: Record<string, string>,
): string {
  let value = translations[key] ?? key;
  if (params) {
    for (const [k, v] of Object.entries(params)) {
      value = value.replaceAll(`{{${k}}}`, v);
    }
  }
  return value;
}
