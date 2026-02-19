import type { ReactNode } from "react";

type AlertVariant = "error" | "success";

interface AlertProps {
  variant: AlertVariant;
  children: ReactNode;
}

const variantClasses: Record<AlertVariant, string> = {
  error:
    "border-destructive/50 bg-destructive/10 text-destructive",
  success:
    "border-green-500/50 bg-green-500/10 text-green-700 dark:text-green-400",
};

export function Alert({ variant, children }: AlertProps) {
  return (
    <div
      className={[
        "rounded-lg border px-3 py-2 text-sm",
        variantClasses[variant],
      ].join(" ")}
    >
      {children}
    </div>
  );
}
