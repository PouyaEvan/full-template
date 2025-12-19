package i18n

import (
	"embed"
	"encoding/json"
	"fmt"
	"sync"
)

//go:embed locales/*.json
var localesFS embed.FS

// Translator handles internationalization of messages
type Translator struct {
	defaultLocale string
	translations  map[string]map[string]string
	mu            sync.RWMutex
}

// NewTranslator creates a new translator with the default locale
func NewTranslator(defaultLocale string) (*Translator, error) {
	t := &Translator{
		defaultLocale: defaultLocale,
		translations:  make(map[string]map[string]string),
	}

	// Load embedded translation files
	if err := t.loadEmbeddedLocales(); err != nil {
		return nil, err
	}

	return t, nil
}

// loadEmbeddedLocales loads all translation files from embedded filesystem
func (t *Translator) loadEmbeddedLocales() error {
	entries, err := localesFS.ReadDir("locales")
	if err != nil {
		// If locales directory doesn't exist, use default translations
		t.loadDefaultTranslations()
		return nil
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		data, err := localesFS.ReadFile("locales/" + entry.Name())
		if err != nil {
			continue
		}

		var translations map[string]string
		if err := json.Unmarshal(data, &translations); err != nil {
			continue
		}

		// Extract locale from filename (e.g., "en.json" -> "en")
		locale := entry.Name()[:len(entry.Name())-5]
		t.translations[locale] = translations
	}

	// If no locales loaded, use defaults
	if len(t.translations) == 0 {
		t.loadDefaultTranslations()
	}

	return nil
}

// loadDefaultTranslations loads hardcoded default translations
func (t *Translator) loadDefaultTranslations() {
	t.translations["en"] = map[string]string{
		"error.invalid_request":          "Invalid request",
		"error.unauthorized":             "Unauthorized access",
		"error.forbidden":                "Access forbidden",
		"error.not_found":                "Resource not found",
		"error.internal":                 "Internal server error",
		"error.invalid_credentials":      "Invalid credentials",
		"error.otp_expired":              "OTP has expired or not found",
		"error.invalid_otp":              "Invalid OTP code",
		"error.user_not_found":           "User not found",
		"error.failed_to_create_user":    "Failed to create user",
		"error.invalid_token":            "Invalid or expired token",
		"error.missing_auth_header":      "Missing authorization header",
		"error.insufficient_permissions": "Insufficient permissions",
		"error.storage_not_configured":   "Storage not configured",
		"error.no_file_uploaded":         "No file uploaded",
		"error.invalid_file_type":        "Invalid file type",
		"error.upload_failed":            "Upload failed",
		"error.payment_failed":           "Payment processing failed",
		"error.2fa_required":             "Two-factor authentication required",
		"error.invalid_2fa_code":         "Invalid 2FA code",
		"error.failed_to_generate_2fa":   "Failed to generate 2FA secret",
		"success.otp_sent":               "OTP sent successfully",
		"success.login":                  "Login successful",
		"success.2fa_enabled":            "Two-factor authentication enabled",
		"success.logout":                 "Logout successful",
		"success.file_uploaded":          "File uploaded successfully",
		"success.payment_processed":      "Payment processed successfully",
	}

	t.translations["fa"] = map[string]string{
		"error.invalid_request":          "درخواست نامعتبر",
		"error.unauthorized":             "دسترسی غیرمجاز",
		"error.forbidden":                "دسترسی ممنوع",
		"error.not_found":                "منبع یافت نشد",
		"error.internal":                 "خطای داخلی سرور",
		"error.invalid_credentials":      "اعتبارنامه نامعتبر",
		"error.otp_expired":              "کد OTP منقضی شده یا یافت نشد",
		"error.invalid_otp":              "کد OTP نامعتبر",
		"error.user_not_found":           "کاربر یافت نشد",
		"error.failed_to_create_user":    "ایجاد کاربر ناموفق بود",
		"error.invalid_token":            "توکن نامعتبر یا منقضی شده",
		"error.missing_auth_header":      "هدر احراز هویت موجود نیست",
		"error.insufficient_permissions": "دسترسی کافی نیست",
		"error.storage_not_configured":   "فضای ذخیره‌سازی پیکربندی نشده",
		"error.no_file_uploaded":         "فایلی آپلود نشده",
		"error.invalid_file_type":        "نوع فایل نامعتبر",
		"error.upload_failed":            "آپلود ناموفق بود",
		"error.payment_failed":           "پردازش پرداخت ناموفق بود",
		"error.2fa_required":             "احراز هویت دو عاملی مورد نیاز است",
		"error.invalid_2fa_code":         "کد 2FA نامعتبر",
		"error.failed_to_generate_2fa":   "تولید کلید 2FA ناموفق بود",
		"success.otp_sent":               "کد OTP با موفقیت ارسال شد",
		"success.login":                  "ورود موفق",
		"success.2fa_enabled":            "احراز هویت دو عاملی فعال شد",
		"success.logout":                 "خروج موفق",
		"success.file_uploaded":          "فایل با موفقیت آپلود شد",
		"success.payment_processed":      "پرداخت با موفقیت انجام شد",
	}
}

// T translates a message key to the specified locale
func (t *Translator) T(locale, key string) string {
	t.mu.RLock()
	defer t.mu.RUnlock()

	// Try to get translation for requested locale
	if translations, ok := t.translations[locale]; ok {
		if msg, ok := translations[key]; ok {
			return msg
		}
	}

	// Fall back to default locale
	if translations, ok := t.translations[t.defaultLocale]; ok {
		if msg, ok := translations[key]; ok {
			return msg
		}
	}

	// Return the key itself if no translation found
	return key
}

// TWithParams translates a message key with parameters
func (t *Translator) TWithParams(locale, key string, params map[string]interface{}) string {
	msg := t.T(locale, key)

	// Simple parameter replacement
	for k, v := range params {
		msg = replaceParam(msg, k, v)
	}

	return msg
}

// replaceParam replaces a parameter placeholder with its value
func replaceParam(msg, key string, value interface{}) string {
	placeholder := fmt.Sprintf("{%s}", key)
	return replaceAll(msg, placeholder, fmt.Sprintf("%v", value))
}

// replaceAll is a simple string replacement function
func replaceAll(s, old, new string) string {
	result := s
	for {
		i := indexOf(result, old)
		if i == -1 {
			break
		}
		result = result[:i] + new + result[i+len(old):]
	}
	return result
}

// indexOf finds the index of a substring
func indexOf(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}

// GetAvailableLocales returns all available locales
func (t *Translator) GetAvailableLocales() []string {
	t.mu.RLock()
	defer t.mu.RUnlock()

	locales := make([]string, 0, len(t.translations))
	for locale := range t.translations {
		locales = append(locales, locale)
	}
	return locales
}

// SetDefaultLocale sets the default fallback locale
func (t *Translator) SetDefaultLocale(locale string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.defaultLocale = locale
}

// AddTranslation adds a translation for a specific locale and key
func (t *Translator) AddTranslation(locale, key, value string) {
	t.mu.Lock()
	defer t.mu.Unlock()

	if _, ok := t.translations[locale]; !ok {
		t.translations[locale] = make(map[string]string)
	}
	t.translations[locale][key] = value
}

// Global translator instance
var globalTranslator *Translator

// InitGlobalTranslator initializes the global translator
func InitGlobalTranslator(defaultLocale string) error {
	var err error
	globalTranslator, err = NewTranslator(defaultLocale)
	return err
}

// Translate translates a key using the global translator
func Translate(locale, key string) string {
	if globalTranslator == nil {
		return key
	}
	return globalTranslator.T(locale, key)
}

// TranslateWithParams translates a key with parameters using the global translator
func TranslateWithParams(locale, key string, params map[string]interface{}) string {
	if globalTranslator == nil {
		return key
	}
	return globalTranslator.TWithParams(locale, key, params)
}
