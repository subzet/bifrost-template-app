import { COUNTRIES } from "../lib/countries";
import { Select } from "../ui/select";

interface CountrySelectProps {
  name: string;
  value?: string;
}

export function CountrySelect({ name, value }: CountrySelectProps) {
  return (
    <Select name={name} defaultValue={value ?? ""}>
      <option value="">—</option>
      {COUNTRIES.map((c) => (
        <option key={c.code} value={c.code}>
          {c.name}
        </option>
      ))}
    </Select>
  );
}
