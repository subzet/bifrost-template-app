import { useEffect, useRef, useState } from "react";
import { Button } from "./button";
import type { ComponentProps } from "react";

type Props = Omit<ComponentProps<typeof Button>, "loading" | "type">;

export function SubmitButton({ children, ...props }: Props) {
  const ref = useRef<HTMLButtonElement>(null);
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    const form = ref.current?.closest("form");
    if (!form) return;
    const handleSubmit = () => setLoading(true);
    form.addEventListener("submit", handleSubmit);
    return () => form.removeEventListener("submit", handleSubmit);
  }, []);

  return (
    <Button ref={ref} type="submit" loading={loading} {...props}>
      {children}
    </Button>
  );
}
