import { Facehash } from "facehash";
import Layout from "./layout";
import { ThemeScript } from "./theme-script";
import { t } from "./i18n";
import { CountrySelect } from "./country-select";

interface EditProfileProps {
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
  error?: string;
  success?: boolean;
  locale: string;
  t: Record<string, string>;
}

export function Head() {
  return (
    <>
      <ThemeScript />
      <title>Edit Profile - MyApp</title>
      <meta name="description" content="Edit your profile" />
    </>
  );
}

const inputClass =
  "flex h-9 w-full rounded-lg border border-input bg-transparent px-2.5 py-1 text-sm placeholder:text-muted-foreground focus-visible:border-ring focus-visible:outline-none focus-visible:ring-[3px] focus-visible:ring-ring/20";

const selectClass =
  "flex h-9 w-full rounded-lg border border-input bg-transparent px-2.5 py-1 text-sm focus-visible:border-ring focus-visible:outline-none focus-visible:ring-[3px] focus-visible:ring-ring/20";

export default function EditProfile({
  user,
  profile,
  error,
  success,
  locale,
  t: translations,
}: EditProfileProps) {
  return (
    <Layout user={user} locale={locale} t={translations}>
      <div className="container flex justify-center py-12">
        <div className="w-full max-w-lg">
          <h1 className="text-2xl font-bold mb-6">{t(translations, "edit.title")}</h1>

          {error && (
            <div className="mb-4 rounded-lg border border-destructive/50 bg-destructive/10 px-3 py-2 text-sm text-destructive">
              {error}
            </div>
          )}

          {success && (
            <div className="mb-4 rounded-lg border border-green-500/50 bg-green-500/10 px-3 py-2 text-sm text-green-700 dark:text-green-400">
              {t(translations, "edit.saved")}
            </div>
          )}

          <form
            method="POST"
            action="/api/user/update"
            encType="multipart/form-data"
            className="space-y-4"
          >
            {/* Avatar */}
            <div className="space-y-2">
              <label className="text-sm font-medium">{t(translations, "edit.avatar")}</label>
              <div className="flex items-center gap-4">
                <div className="w-16 h-16 rounded-full overflow-hidden flex-shrink-0 border border-border">
                  {profile.avatarURL ? (
                    <img
                      src={profile.avatarURL}
                      alt="avatar"
                      className="w-full h-full object-cover"
                    />
                  ) : (
                    <Facehash name={profile.email} size={64} />
                  )}
                </div>
                <input
                  type="file"
                  name="avatar"
                  accept="image/jpeg,image/png,image/gif,image/webp"
                  className="text-sm text-muted-foreground file:mr-3 file:rounded-md file:border file:border-border file:bg-transparent file:px-3 file:py-1 file:text-sm file:font-medium"
                />
              </div>
            </div>

            {/* Handle */}
            <div className="space-y-1.5">
              <label htmlFor="handle" className="text-sm font-medium">
                {t(translations, "edit.handle")}
              </label>
              <input
                id="handle"
                type="text"
                name="handle"
                defaultValue={profile.handle}
                className={inputClass}
                required
                pattern="[a-z0-9][a-z0-9_\-]{2,29}"
              />
            </div>

            {/* Display name */}
            <div className="space-y-1.5">
              <label htmlFor="display_name" className="text-sm font-medium">
                {t(translations, "edit.displayName")}
              </label>
              <input
                id="display_name"
                type="text"
                name="display_name"
                defaultValue={profile.displayName}
                className={inputClass}
              />
            </div>

            {/* Bio */}
            <div className="space-y-1.5">
              <label htmlFor="bio" className="text-sm font-medium">
                {t(translations, "edit.bio")}
              </label>
              <textarea
                id="bio"
                name="bio"
                defaultValue={profile.bio}
                rows={3}
                className="flex w-full rounded-lg border border-input bg-transparent px-2.5 py-1.5 text-sm placeholder:text-muted-foreground focus-visible:border-ring focus-visible:outline-none focus-visible:ring-[3px] focus-visible:ring-ring/20"
              />
            </div>

            {/* Country */}
            <div className="space-y-1.5">
              <label htmlFor="country" className="text-sm font-medium">
                {t(translations, "edit.country")}
              </label>
              <CountrySelect name="country" value={profile.country} className={selectClass} />
            </div>

            {/* Social links */}
            <div className="space-y-2">
              <h3 className="text-sm font-medium">{t(translations, "edit.socials")}</h3>
              <input
                type="text"
                name="instagram"
                defaultValue={profile.socialLinks.instagram}
                placeholder={t(translations, "edit.instagram")}
                className={inputClass}
              />
              <input
                type="text"
                name="facebook"
                defaultValue={profile.socialLinks.facebook}
                placeholder={t(translations, "edit.facebook")}
                className={inputClass}
              />
              <input
                type="text"
                name="linkedin"
                defaultValue={profile.socialLinks.linkedin}
                placeholder={t(translations, "edit.linkedin")}
                className={inputClass}
              />
              <input
                type="text"
                name="x"
                defaultValue={profile.socialLinks.x}
                placeholder={t(translations, "edit.x")}
                className={inputClass}
              />
            </div>

            <button
              type="submit"
              className="inline-flex h-10 w-full items-center justify-center rounded-lg bg-primary px-2.5 text-sm font-medium text-primary-foreground"
            >
              {t(translations, "edit.submit")}
            </button>
          </form>
        </div>
      </div>
    </Layout>
  );
}
