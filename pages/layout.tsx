import type { ReactNode } from "react";
import { ThemeToggle } from "./theme-toggle";
import { t } from "./i18n";
import "./app.css";

interface LayoutProps {
  user?: { email: string };
  locale: string;
  t: Record<string, string>;
  children: ReactNode;
}

export default function Layout({ user, locale, t: translations, children }: LayoutProps) {
  return (
    <div className="min-h-screen flex flex-col">
      <nav className="container flex justify-between py-4">
        <div className="flex items-center gap-1">
          <a href="/" className="mr-2 font-bold text-foreground">
            MyApp
          </a>
          <ThemeToggle />
        </div>
        <div className="flex items-center gap-1">
          <form method="POST" action="/api/set-lang">
            <input type="hidden" name="lang" value={locale === "es" ? "en" : "es"} />
            <button
              type="submit"
              className="inline-flex items-center justify-center rounded-lg px-2 py-1.5 text-sm text-muted-foreground underline-offset-4 hover:underline"
            >
              {locale === "es" ? "EN" : "ES"}
            </button>
          </form>
          {user ? (
            <>
              <span className="text-sm text-muted-foreground">{user.email}</span>
              <form method="POST" action="/api/logout">
                <button
                  type="submit"
                  className="inline-flex items-center justify-center rounded-lg px-2 py-1.5 text-sm text-muted-foreground underline-offset-4 hover:underline"
                >
                  {t(translations, "nav.logout")}
                </button>
              </form>
            </>
          ) : (
            <>
              <a
                href="/login"
                className="inline-flex items-center justify-center rounded-lg px-2 py-1.5 text-sm text-muted-foreground underline-offset-4 hover:underline"
              >
                {t(translations, "nav.login")}
              </a>
              <a
                href="/signup"
                className="inline-flex items-center justify-center rounded-lg bg-primary text-primary-foreground px-2.5 py-1.5 text-sm font-medium"
              >
                {t(translations, "nav.signup")}
              </a>
            </>
          )}
        </div>
      </nav>

      <main className="flex-1">{children}</main>

      <footer className="container flex flex-1 items-end py-6">
        <div className="flex w-full items-center justify-between">
          <span className="text-sm text-muted-foreground">
            {t(translations, "footer.builtWith")}
          </span>
          <a href="/" className="text-sm text-muted-foreground underline-offset-4 hover:underline">
            &copy; {new Date().getFullYear()} MyApp.
          </a>
        </div>
      </footer>
    </div>
  );
}
