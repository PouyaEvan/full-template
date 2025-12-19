package i18n

import (
	"strings"

	"github.com/gofiber/fiber/v2"
)

// Config holds the configuration for the i18n middleware
type Config struct {
	// DefaultLocale is the fallback locale
	DefaultLocale string
	// SupportedLocales is a list of supported locale codes
	SupportedLocales []string
	// LocaleContextKey is the key used to store locale in fiber context
	LocaleContextKey string
}

// DefaultConfig returns the default configuration
func DefaultConfig() Config {
	return Config{
		DefaultLocale:    "en",
		SupportedLocales: []string{"en", "fa"},
		LocaleContextKey: "locale",
	}
}

// Middleware creates a Fiber middleware that detects and sets the locale
func Middleware(config ...Config) fiber.Handler {
	cfg := DefaultConfig()
	if len(config) > 0 {
		cfg = config[0]
	}

	return func(c *fiber.Ctx) error {
		locale := detectLocale(c, cfg)
		c.Locals(cfg.LocaleContextKey, locale)
		return c.Next()
	}
}

// detectLocale detects the locale from the request
func detectLocale(c *fiber.Ctx, cfg Config) string {
	// 1. Check query parameter (?lang=fa)
	if lang := c.Query("lang"); lang != "" {
		if isSupported(lang, cfg.SupportedLocales) {
			return lang
		}
	}

	// 2. Check X-Locale header
	if lang := c.Get("X-Locale"); lang != "" {
		if isSupported(lang, cfg.SupportedLocales) {
			return lang
		}
	}

	// 3. Check Accept-Language header
	acceptLang := c.Get("Accept-Language")
	if acceptLang != "" {
		lang := parseAcceptLanguage(acceptLang, cfg.SupportedLocales)
		if lang != "" {
			return lang
		}
	}

	return cfg.DefaultLocale
}

// parseAcceptLanguage parses the Accept-Language header and returns the best match
func parseAcceptLanguage(header string, supported []string) string {
	// Simple parsing: split by comma and check each language
	langs := strings.Split(header, ",")
	for _, lang := range langs {
		// Remove quality factor (e.g., "en-US;q=0.9" -> "en-US")
		lang = strings.TrimSpace(strings.Split(lang, ";")[0])

		// Try exact match first
		if isSupported(lang, supported) {
			return lang
		}

		// Try language code without region (e.g., "en-US" -> "en")
		if idx := strings.Index(lang, "-"); idx != -1 {
			langCode := lang[:idx]
			if isSupported(langCode, supported) {
				return langCode
			}
		}
	}

	return ""
}

// isSupported checks if a locale is in the supported list
func isSupported(locale string, supported []string) bool {
	for _, s := range supported {
		if strings.EqualFold(locale, s) {
			return true
		}
	}
	return false
}

// GetLocale gets the locale from fiber context
func GetLocale(c *fiber.Ctx) string {
	locale, ok := c.Locals("locale").(string)
	if !ok {
		return "en"
	}
	return locale
}

// LocalizedError returns a fiber error response with translated message
func LocalizedError(c *fiber.Ctx, status int, key string) error {
	locale := GetLocale(c)
	return c.Status(status).JSON(fiber.Map{
		"error": Translate(locale, key),
	})
}

// LocalizedSuccess returns a fiber success response with translated message
func LocalizedSuccess(c *fiber.Ctx, key string, data ...fiber.Map) error {
	locale := GetLocale(c)
	response := fiber.Map{
		"message": Translate(locale, key),
	}
	if len(data) > 0 {
		for k, v := range data[0] {
			response[k] = v
		}
	}
	return c.JSON(response)
}
