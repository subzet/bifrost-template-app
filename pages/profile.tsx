import { Facehash } from "facehash";
import Layout from "./layout";
import { ThemeScript } from "./theme-script";
import { t } from "./i18n";
import { countryName } from "./countries";

interface ProfileProps {
  user?: { email: string; handle: string };
  profile: {
    handle: string;
    displayName: string;
    bio: string;
    country: string;
    email: string;
    avatarURL: string;
    socialLinks: { instagram: string; facebook: string; linkedin: string; x: string };
  };
  isOwner: boolean;
  locale: string;
  t: Record<string, string>;
}

export function Head() {
  return (
    <>
      <ThemeScript />
      <title>Profile - MyApp</title>
      <meta name="description" content="User profile" />
    </>
  );
}

export default function Profile({ user, profile, isOwner, locale, t: translations }: ProfileProps) {
  const hasInfo =
    profile.bio ||
    profile.country ||
    profile.socialLinks.instagram ||
    profile.socialLinks.facebook ||
    profile.socialLinks.linkedin ||
    profile.socialLinks.x;

  return (
    <Layout user={user} locale={locale} t={translations}>
      <div className="container py-12 max-w-2xl mx-auto">
        <div className="flex items-start justify-between mb-8">
          <div className="flex items-center gap-4">
            <div className="w-20 h-20 rounded-full overflow-hidden flex-shrink-0">
              {profile.avatarURL ? (
                <img
                  src={profile.avatarURL}
                  alt={`@${profile.handle}`}
                  className="w-full h-full object-cover"
                />
              ) : (
                <Facehash name={profile.email} size={80} />
              )}
            </div>
            <div>
              <h1 className="text-2xl font-bold">
                {profile.displayName || `@${profile.handle}`}
              </h1>
              {profile.displayName && (
                <p className="text-muted-foreground">@{profile.handle}</p>
              )}
            </div>
          </div>
          {isOwner && (
            <a
              href={`/user/${profile.handle}/edit`}
              className="inline-flex items-center justify-center rounded-lg border border-border px-4 py-2 text-sm font-medium hover:bg-muted"
            >
              {t(translations, "profile.editButton")}
            </a>
          )}
        </div>

        {!hasInfo && (
          <p className="text-muted-foreground">{t(translations, "profile.noInfo")}</p>
        )}

        {profile.bio && (
          <div className="mb-6">
            <h2 className="text-sm font-medium text-muted-foreground mb-1">
              {t(translations, "profile.bio")}
            </h2>
            <p>{profile.bio}</p>
          </div>
        )}

        {profile.country && (
          <div className="mb-6">
            <h2 className="text-sm font-medium text-muted-foreground mb-1">
              {t(translations, "profile.country")}
            </h2>
            <p>{countryName(profile.country)}</p>
          </div>
        )}

        {(profile.socialLinks.instagram ||
          profile.socialLinks.facebook ||
          profile.socialLinks.linkedin ||
          profile.socialLinks.x) && (
          <div className="mb-6">
            <h2 className="text-sm font-medium text-muted-foreground mb-2">
              {t(translations, "profile.socials")}
            </h2>
            <div className="flex flex-wrap gap-3">
              {profile.socialLinks.instagram && (
                <a
                  href={profile.socialLinks.instagram}
                  className="text-sm underline underline-offset-4"
                  target="_blank"
                  rel="noopener noreferrer"
                >
                  Instagram
                </a>
              )}
              {profile.socialLinks.facebook && (
                <a
                  href={profile.socialLinks.facebook}
                  className="text-sm underline underline-offset-4"
                  target="_blank"
                  rel="noopener noreferrer"
                >
                  Facebook
                </a>
              )}
              {profile.socialLinks.linkedin && (
                <a
                  href={profile.socialLinks.linkedin}
                  className="text-sm underline underline-offset-4"
                  target="_blank"
                  rel="noopener noreferrer"
                >
                  LinkedIn
                </a>
              )}
              {profile.socialLinks.x && (
                <a
                  href={profile.socialLinks.x}
                  className="text-sm underline underline-offset-4"
                  target="_blank"
                  rel="noopener noreferrer"
                >
                  X
                </a>
              )}
            </div>
          </div>
        )}
      </div>
    </Layout>
  );
}
