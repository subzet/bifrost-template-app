import Layout from "./layout";
import { ThemeScript } from "./theme-script";
import { t } from "./lib/i18n";
import { Alert } from "./ui/alert";
import { SubmitButton } from "./ui/submit-button";
import { Card } from "./ui/card";
import { FormField } from "./ui/form-field";
import { Input } from "./ui/input";

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
    <Layout user={user} locale={locale} t={translations} hideAuthLinks>
      <div className="container flex justify-center py-24">
        <Card className="w-full max-w-sm">
          <h2 className="text-center text-lg font-medium">
            {t(translations, "signup.title")}
          </h2>

          {error && (
            <div className="mt-4">
              <Alert variant="error">{error}</Alert>
            </div>
          )}

          <form method="POST" action="/api/signup" className="mt-6 space-y-4">
            <FormField label={t(translations, "signup.handle")} htmlFor="handle">
              <Input
                id="handle"
                type="text"
                name="handle"
                placeholder="yourhandle"
                required
                pattern="[a-z0-9][a-z0-9_\-]{2,29}"
              />
            </FormField>

            <FormField label={t(translations, "signup.email")} htmlFor="email">
              <Input
                id="email"
                type="email"
                name="email"
                placeholder="you@example.com"
                required
              />
            </FormField>

            <FormField label={t(translations, "signup.password")} htmlFor="password">
              <Input
                id="password"
                type="password"
                name="password"
                placeholder="••••••••"
                required
                minLength={8}
              />
            </FormField>

            <FormField label={t(translations, "signup.confirmPassword")} htmlFor="confirm_password">
              <Input
                id="confirm_password"
                type="password"
                name="confirm_password"
                placeholder="••••••••"
                required
                minLength={8}
              />
            </FormField>

            <SubmitButton fullWidth>
              {t(translations, "signup.submit")}
            </SubmitButton>
          </form>

          <p className="mt-4 text-center text-sm text-muted-foreground">
            {t(translations, "signup.hasAccount")}{" "}
            <a href="/login" className="text-foreground underline underline-offset-4 hover:text-foreground/80">
              {t(translations, "signup.loginLink")}
            </a>
          </p>
        </Card>
      </div>
    </Layout>
  );
}
