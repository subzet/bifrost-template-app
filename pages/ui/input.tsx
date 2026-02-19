import type { ComponentProps } from "react";

const inputClass =
  "flex h-9 w-full rounded-lg border border-input bg-transparent px-2.5 py-1 text-sm placeholder:text-muted-foreground focus-visible:border-ring focus-visible:outline-none focus-visible:ring-[3px] focus-visible:ring-ring/20";

export function Input({ className, ...props }: ComponentProps<"input">) {
  return (
    <input
      className={[inputClass, className].filter(Boolean).join(" ")}
      {...props}
    />
  );
}
