import type { ComponentProps } from "react";

type Variant = "primary" | "outline" | "ghost";

interface ButtonProps extends ComponentProps<"button"> {
  variant?: Variant;
  size?: "default" | "sm";
  fullWidth?: boolean;
  loading?: boolean;
}

const base = "inline-flex items-center justify-center rounded-lg text-sm font-medium transition-opacity";

const variantClasses: Record<Variant, string> = {
  primary: "bg-primary text-primary-foreground",
  outline: "border border-border hover:bg-muted",
  ghost: "text-muted-foreground underline-offset-4 hover:underline",
};

const sizeClasses: Record<Variant, Record<"default" | "sm", string>> = {
  primary: { default: "h-10 px-2.5", sm: "px-2.5 py-1.5" },
  outline: { default: "px-4 py-2", sm: "px-3 py-1.5" },
  ghost: { default: "px-2 py-1.5", sm: "px-2 py-1" },
};

function Spinner() {
  return (
    <svg
      className="animate-spin -ml-1 mr-2 h-4 w-4 shrink-0"
      xmlns="http://www.w3.org/2000/svg"
      fill="none"
      viewBox="0 0 24 24"
      aria-hidden="true"
    >
      <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4" />
      <path
        className="opacity-75"
        fill="currentColor"
        d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
      />
    </svg>
  );
}

export function Button({
  variant = "primary",
  size = "default",
  fullWidth,
  loading,
  disabled,
  className,
  children,
  ...props
}: ButtonProps) {
  return (
    <button
      disabled={disabled || loading}
      className={[
        base,
        variantClasses[variant],
        sizeClasses[variant][size],
        fullWidth && "w-full",
        (disabled || loading) && "opacity-60 cursor-not-allowed",
        className,
      ]
        .filter(Boolean)
        .join(" ")}
      {...props}
    >
      {loading && <Spinner />}
      {children}
    </button>
  );
}

export function buttonClass(
  variant: Variant = "primary",
  size: "default" | "sm" = "default",
  fullWidth = false,
): string {
  return [
    base,
    variantClasses[variant],
    sizeClasses[variant][size],
    fullWidth && "w-full",
  ]
    .filter(Boolean)
    .join(" ");
}
