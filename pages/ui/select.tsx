import type { ComponentProps } from "react";

export function Select({ className, ...props }: ComponentProps<"select">) {
  return (
    <select
      className={[
        "flex h-9 w-full rounded-lg border border-input bg-transparent px-2.5 py-1 text-sm focus-visible:border-ring focus-visible:outline-none focus-visible:ring-[3px] focus-visible:ring-ring/20",
        className,
      ]
        .filter(Boolean)
        .join(" ")}
      {...props}
    />
  );
}
