'use client';

import { useTransition } from 'react';
import { useRouter } from 'next/navigation';
import { locales, type Locale } from '@/i18n';

export function LanguageSwitcher() {
  const [isPending, startTransition] = useTransition();
  const router = useRouter();

  const setLocale = (locale: Locale) => {
    startTransition(() => {
      // Set cookie for server-side locale detection
      document.cookie = `locale=${locale}; path=/; max-age=31536000`;
      // Use router.refresh() to re-render the page with new locale without full reload
      router.refresh();
    });
  };

  return (
    <div className="flex items-center gap-2">
      {locales.map((locale) => (
        <button
          key={locale}
          onClick={() => setLocale(locale)}
          disabled={isPending}
          className="px-2 py-1 text-sm rounded-md hover:bg-gray-100 dark:hover:bg-gray-800 disabled:opacity-50"
        >
          {locale === 'en' ? 'English' : 'فارسی'}
        </button>
      ))}
    </div>
  );
}
