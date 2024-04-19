// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// License: BSD 3-Clause License

package translation

import (
	"fmt"
	"os"

	"github.com/bytedance/sonic"
)

// Translations is a map of language codes to their respective translation maps.
var Translations map[string]map[string]string

// LoadTranslations loads translations from a single JSON file.
func LoadTranslations(filePath string) error {
	// Note: This method is better because on Windows, long paths are not allowed by default, unlike on Unix/Linux.
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
func Translate(lang, key string, args ...interface{}) string {
	translation, exists := Translations[lang][key]
	if !exists {
		return key // Fallback to the key itself if translation does not exist
	}
	return fmt.Sprintf(translation, args...)
}
