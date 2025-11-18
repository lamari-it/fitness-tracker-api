package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
	"sync"
)

type I18n struct {
	defaultLanguage string
	translations    map[string]map[string]interface{}
	mu              sync.RWMutex
}

var i18nInstance *I18n
var i18nOnce sync.Once

// GetI18n returns the singleton instance of I18n
func GetI18n() *I18n {
	i18nOnce.Do(func() {
		i18nInstance = &I18n{
			defaultLanguage: "en",
			translations:    make(map[string]map[string]interface{}),
		}
		i18nInstance.loadTranslations()
	})
	return i18nInstance
}

// loadTranslations loads all translation files from the locales directory
func (i *I18n) loadTranslations() {
	supportedLanguages := []string{"en", "es", "fr"}

	for _, lang := range supportedLanguages {
		filePath := filepath.Join("locales", lang, "messages.json")
		data, err := ioutil.ReadFile(filePath)
		if err != nil {
			fmt.Printf("Error reading translation file %s: %v\n", filePath, err)
			continue
		}

		var translations map[string]interface{}
		if err := json.Unmarshal(data, &translations); err != nil {
			fmt.Printf("Error parsing translation file %s: %v\n", filePath, err)
			continue
		}

		i.mu.Lock()
		i.translations[lang] = translations
		i.mu.Unlock()
	}
}

// GetSupportedLanguages returns list of supported languages
func (i *I18n) GetSupportedLanguages() []string {
	i.mu.RLock()
	defer i.mu.RUnlock()

	languages := make([]string, 0, len(i.translations))
	for lang := range i.translations {
		languages = append(languages, lang)
	}
	return languages
}

// IsLanguageSupported checks if a language is supported
func (i *I18n) IsLanguageSupported(lang string) bool {
	i.mu.RLock()
	defer i.mu.RUnlock()

	_, exists := i.translations[lang]
	return exists
}

// Translate translates a message key to the specified language
func (i *I18n) Translate(lang, key string, args ...interface{}) string {
	i.mu.RLock()
	defer i.mu.RUnlock()

	// If language not supported, use default
	if !i.IsLanguageSupported(lang) {
		lang = i.defaultLanguage
	}

	// Get translation for the language
	langTranslations, exists := i.translations[lang]
	if !exists {
		return key // Return key if no translations found
	}

	// Navigate through nested keys (e.g., "auth.login_success")
	keys := strings.Split(key, ".")
	current := langTranslations

	for _, k := range keys {
		if nested, ok := current[k]; ok {
			if nestedMap, isMap := nested.(map[string]interface{}); isMap {
				current = nestedMap
			} else if str, isString := nested.(string); isString {
				// Format with arguments if provided
				if len(args) > 0 {
					return fmt.Sprintf(str, args...)
				}
				return str
			}
		} else {
			// Key not found, try fallback to default language
			if lang != i.defaultLanguage {
				return i.Translate(i.defaultLanguage, key, args...)
			}
			return key // Return key if not found in any language
		}
	}

	return key // Return key if we couldn't find the translation
}

// T is a shorthand for Translate
func (i *I18n) T(lang, key string, args ...interface{}) string {
	return i.Translate(lang, key, args...)
}

// GetDefaultLanguage returns the default language
func (i *I18n) GetDefaultLanguage() string {
	return i.defaultLanguage
}

// SetDefaultLanguage sets the default language
func (i *I18n) SetDefaultLanguage(lang string) {
	i.mu.Lock()
	defer i.mu.Unlock()
	i.defaultLanguage = lang
}

// ReloadTranslations reloads all translation files
func (i *I18n) ReloadTranslations() {
	i.mu.Lock()
	defer i.mu.Unlock()
	i.translations = make(map[string]map[string]interface{})
	i.loadTranslations()
}
