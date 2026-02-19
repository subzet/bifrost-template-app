import type { ReactNode } from "react";
import { ThemeToggle } from "./theme-toggle";
import { t } from "./lib/i18n";
import { Button, buttonClass } from "./ui/button";
import "./app.css";

interface LayoutProps {
  user?: { email: string; handle: string };
  locale: string;
  t: Record<string, string>;
  children: ReactNode;
  hideAuthLinks?: boolean;
}

export default function Layout({ user, locale, t: translations, children, hideAuthLinks }: LayoutProps) {
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
            <Button variant="ghost" size="sm" type="submit">
              {locale === "es" ? "EN" : "ES"}
            </Button>
          </form>
          {user ? (
            <>
              <a href={`/user/${user.handle}`} className="text-sm text-muted-foreground underline-offset-4 hover:underline">@{user.handle}</a>
              <form method="POST" action="/api/logout">
                <Button variant="ghost" size="sm" type="submit">
                  {t(translations, "nav.logout")}
                </Button>
              </form>
            </>
          ) : !hideAuthLinks && (
            <>
              <a href="/login" className={buttonClass("ghost", "sm")}>
                {t(translations, "nav.login")}
              </a>
              <a href="/signup" className={buttonClass("primary", "sm")}>
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
