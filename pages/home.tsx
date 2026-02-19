import Layout from "./layout";
import { ThemeScript } from "./theme-script";
import { t } from "./lib/i18n";

interface HomeProps {
  user?: { email: string; handle: string };
  locale: string;
  t: Record<string, string>;
}

export function Head() {
  return (
    <>
      <ThemeScript />
      <title>MyApp</title>
      <meta name="description" content="Welcome to MyApp" />
    </>
  );
}

export default function Home({ user, locale, t: translations }: HomeProps) {
  return (
    <Layout user={user} locale={locale} t={translations}>
      <div className="container flex flex-col items-center justify-center py-24 text-center">
        <h1 className="text-4xl font-bold tracking-tight">
          {t(translations, "home.title")}
        </h1>
        {user ? (
          <p className="mt-4 text-lg text-muted-foreground">
            {t(translations, "home.greeting", { email: user.email })}
          </p>
        ) : (
          <p className="mt-4 text-lg text-muted-foreground">
            {t(translations, "home.cta")}
          </p>
        )}
      </div>
    </Layout>
  );
}
