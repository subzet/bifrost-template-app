import Layout from "./layout";
import { ThemeScript } from "./theme-script";
import { t } from "./i18n";

interface SignupProps {
  user?: { email: string; handle: string };
  error?: string;
  locale: string;
  t: Record<string, string>;
}

export function Head() {
  return (
    <>
      <ThemeScript />
      <title>Sign Up - MyApp</title>
      <meta name="description" content="Create a MyApp account" />
    </>
  );
}

export default function Signup({ user, error, locale, t: translations }: SignupProps) {
  return (
    <Layout user={user} locale={locale} t={translations}>
      <div className="container flex justify-center py-24">
        <div className="w-full max-w-sm rounded-xl border border-border bg-card p-6 text-card-foreground">
          <h2 className="text-center text-lg font-medium">
            {t(translations, "signup.title")}
          </h2>

          {error && (
            <div className="mt-4 rounded-lg border border-destructive/50 bg-destructive/10 px-3 py-2 text-sm text-destructive">
              {error}
            </div>
          )}

          <form method="POST" action="/api/signup" className="mt-6 space-y-4">
            <div className="space-y-1.5">
              <label htmlFor="handle" className="text-sm font-medium">
                {t(translations, "signup.handle")}
              </label>
              <input
                id="handle"
                type="text"
                name="handle"
                placeholder="yourhandle"
                className="flex h-9 w-full rounded-lg border border-input bg-transparent px-2.5 py-1 text-sm placeholder:text-muted-foreground focus-visible:border-ring focus-visible:outline-none focus-visible:ring-[3px] focus-visible:ring-ring/20"
                required
                pattern="[a-z0-9][a-z0-9_\-]{2,29}"
              />
            </div>

            <div className="space-y-1.5">
              <label htmlFor="email" className="text-sm font-medium">
                {t(translations, "signup.email")}
              </label>
              <input
                id="email"
                type="email"
                name="email"
                placeholder="you@example.com"
                className="flex h-9 w-full rounded-lg border border-input bg-transparent px-2.5 py-1 text-sm placeholder:text-muted-foreground focus-visible:border-ring focus-visible:outline-none focus-visible:ring-[3px] focus-visible:ring-ring/20"
                required
              />
            </div>

            <div className="space-y-1.5">
              <label htmlFor="password" className="text-sm font-medium">
                {t(translations, "signup.password")}
              </label>
              <input
                id="password"
                type="password"
                name="password"
                placeholder="••••••••"
                className="flex h-9 w-full rounded-lg border border-input bg-transparent px-2.5 py-1 text-sm placeholder:text-muted-foreground focus-visible:border-ring focus-visible:outline-none focus-visible:ring-[3px] focus-visible:ring-ring/20"
                required
                minLength={8}
              />
            </div>

            <div className="space-y-1.5">
              <label htmlFor="confirm_password" className="text-sm font-medium">
                {t(translations, "signup.confirmPassword")}
              </label>
              <input
                id="confirm_password"
                type="password"
                name="confirm_password"
                placeholder="••••••••"
                className="flex h-9 w-full rounded-lg border border-input bg-transparent px-2.5 py-1 text-sm placeholder:text-muted-foreground focus-visible:border-ring focus-visible:outline-none focus-visible:ring-[3px] focus-visible:ring-ring/20"
                required
                minLength={8}
              />
            </div>

            <button
              type="submit"
              className="inline-flex h-10 w-full items-center justify-center rounded-lg bg-primary px-2.5 text-sm font-medium text-primary-foreground"
            >
              {t(translations, "signup.submit")}
            </button>
          </form>

          <p className="mt-4 text-center text-sm text-muted-foreground">
            {t(translations, "signup.hasAccount")}{" "}
            <a href="/login" className="text-foreground underline underline-offset-4 hover:text-foreground/80">
              {t(translations, "signup.loginLink")}
            </a>
          </p>
        </div>
      </div>
    </Layout>
  );
}
