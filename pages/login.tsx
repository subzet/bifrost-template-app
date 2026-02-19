import Layout from "./layout";
import { ThemeScript } from "./theme-script";
import { t } from "./lib/i18n";
import { Alert } from "./ui/alert";
import { SubmitButton } from "./ui/submit-button";
import { Card } from "./ui/card";
import { FormField } from "./ui/form-field";
import { Input } from "./ui/input";

interface LoginProps {
  user?: { email: string; handle: string };
  error?: string;
  locale: string;
  t: Record<string, string>;
}

export function Head() {
  return (
    <>
      <ThemeScript />
      <title>Login - MyApp</title>
      <meta name="description" content="Login to MyApp" />
    </>
  );
}

export default function Login({ user, error, locale, t: translations }: LoginProps) {
  return (
    <Layout user={user} locale={locale} t={translations} hideAuthLinks>
      <div className="container flex justify-center py-24">
        <Card className="w-full max-w-sm">
          <h2 className="text-center text-lg font-medium">
            {t(translations, "login.title")}
          </h2>

          {error && (
            <div className="mt-4">
              <Alert variant="error">{error}</Alert>
            </div>
          )}

          <form method="POST" action="/api/login" className="mt-6 space-y-4">
            <FormField label={t(translations, "login.email")} htmlFor="email">
              <Input
                id="email"
                type="email"
                name="email"
                placeholder="you@example.com"
                required
              />
            </FormField>

            <FormField label={t(translations, "login.password")} htmlFor="password">
              <Input
                id="password"
                type="password"
                name="password"
                placeholder="••••••••"
                required
              />
            </FormField>

            <SubmitButton fullWidth>
              {t(translations, "login.submit")}
            </SubmitButton>
          </form>

          <p className="mt-4 text-center text-sm text-muted-foreground">
            {t(translations, "login.noAccount")}{" "}
            <a href="/signup" className="text-foreground underline underline-offset-4 hover:text-foreground/80">
              {t(translations, "login.signupLink")}
            </a>
          </p>
        </Card>
      </div>
    </Layout>
  );
}
