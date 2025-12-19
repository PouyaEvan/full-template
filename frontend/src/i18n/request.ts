import { getRequestConfig } from 'next-intl/server';
import { cookies, headers } from 'next/headers';

export const locales = ['en', 'fa'] as const;
export type Locale = (typeof locales)[number];
export const defaultLocale: Locale = 'en';

export default getRequestConfig(async () => {
  // Try to get locale from cookie
  let locale: Locale = defaultLocale;

  const cookieStore = cookies();
  const localeCookie = cookieStore.get('locale')?.value;

  if (localeCookie && locales.includes(localeCookie as Locale)) {
    locale = localeCookie as Locale;
  } else {
    // Try to get from Accept-Language header
    const headersList = headers();
    const acceptLanguage = headersList.get('accept-language');

    if (acceptLanguage) {
      const preferredLocale = acceptLanguage
        .split(',')
        .map((lang) => lang.split(';')[0].trim().split('-')[0])
        .find((lang) => locales.includes(lang as Locale));

      if (preferredLocale) {
        locale = preferredLocale as Locale;
      }
    }
  }

  return {
    locale,
    messages: (await import(`../../messages/${locale}.json`)).default
  };
});
