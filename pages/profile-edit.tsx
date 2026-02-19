import { Facehash } from "facehash";
import Layout from "./layout";
import { ThemeScript } from "./theme-script";
import { t } from "./lib/i18n";
import { CountrySelect } from "./components/country-select";
import { Alert } from "./ui/alert";
import { SubmitButton } from "./ui/submit-button";
import { FormField } from "./ui/form-field";
import { Input } from "./ui/input";
import { Textarea } from "./ui/textarea";

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
            <div className="mb-4">
              <Alert variant="error">{error}</Alert>
            </div>
          )}

          {success && (
            <div className="mb-4">
              <Alert variant="success">{t(translations, "edit.saved")}</Alert>
            </div>
          )}

          <form
            method="POST"
            action="/api/user/update"
            encType="multipart/form-data"
            className="space-y-4"
          >
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

            <FormField label={t(translations, "edit.handle")} htmlFor="handle">
              <Input
                id="handle"
                type="text"
                name="handle"
                defaultValue={profile.handle}
                required
                pattern="[a-z0-9][a-z0-9_\-]{2,29}"
              />
            </FormField>

            <FormField label={t(translations, "edit.displayName")} htmlFor="display_name">
              <Input
                id="display_name"
                type="text"
                name="display_name"
                defaultValue={profile.displayName}
              />
            </FormField>

            <FormField label={t(translations, "edit.bio")} htmlFor="bio">
              <Textarea
                id="bio"
                name="bio"
                defaultValue={profile.bio}
                rows={3}
              />
            </FormField>

            <FormField label={t(translations, "edit.country")} htmlFor="country">
              <CountrySelect name="country" value={profile.country} />
            </FormField>

            <div className="space-y-2">
              <h3 className="text-sm font-medium">{t(translations, "edit.socials")}</h3>
              <Input
                type="text"
                name="instagram"
                defaultValue={profile.socialLinks.instagram}
                placeholder={t(translations, "edit.instagram")}
              />
              <Input
                type="text"
                name="facebook"
                defaultValue={profile.socialLinks.facebook}
                placeholder={t(translations, "edit.facebook")}
              />
              <Input
                type="text"
                name="linkedin"
                defaultValue={profile.socialLinks.linkedin}
                placeholder={t(translations, "edit.linkedin")}
              />
              <Input
                type="text"
                name="x"
                defaultValue={profile.socialLinks.x}
                placeholder={t(translations, "edit.x")}
              />
            </div>

            <SubmitButton fullWidth>
              {t(translations, "edit.submit")}
            </SubmitButton>
          </form>
        </div>
      </div>
    </Layout>
  );
}
