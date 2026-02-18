import { COUNTRIES } from "./countries";

interface CountrySelectProps {
  name: string;
  value?: string;
  className?: string;
}

export function CountrySelect({ name, value, className }: CountrySelectProps) {
  return (
    <select name={name} defaultValue={value ?? ""} className={className}>
      <option value="">—</option>
      {COUNTRIES.map((c) => (
        <option key={c.code} value={c.code}>
          {c.name}
        </option>
      ))}
    </select>
  );
}
