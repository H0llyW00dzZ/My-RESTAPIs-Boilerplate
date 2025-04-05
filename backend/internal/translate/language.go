// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// By accessing or using this software, you agree to be bound by the terms
// of the License Agreement, which you can find at LICENSE files.

package translation

import (
	"fmt"
	"os"

	"github.com/bytedance/sonic"
)

// Translations is a map of language codes to their respective translation maps.
//
// Note: There is a performance cost to using this translation mechanism, as it can grow easily if there is a lot of data (approximately 1MB or more).
var Translations map[string]map[string]string

// LoadTranslations loads translations from a single JSON file.
func LoadTranslations(filePath string) error {
	// Note: This method is better because on Windows, long paths are not allowed by default, unlike on Unix/Linux.
	// Also note that Ignore false positives reported by code scanners (e.g., CodeQL or other scanner tools) that are not 100% accurate.
	// For example, got detected "G304 (CWE-22): Potential file inclusion via variable".
	content, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	err = sonic.Unmarshal(content, &Translations)
	if err != nil {
		return err
	}

	return nil
}

// Translate returns the translated string for the given key and language.
func Translate(lang, key string, args ...any) string {
	translation, exists := Translations[lang][key]
	if !exists {
		return key // Fallback to the key itself if translation does not exist
	}
	return fmt.Sprintf(translation, args...)
}
