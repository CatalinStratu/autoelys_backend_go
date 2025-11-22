package utils

import (
	"regexp"
	"strings"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

// GenerateSlug creates a URL-friendly slug from a string
func GenerateSlug(text string) string {
	// Convert to lowercase
	slug := strings.ToLower(text)

	// Remove diacritics (ă, â, î, ș, ț, etc.)
	slug = removeDiacritics(slug)

	// Replace spaces and special characters with hyphens
	reg := regexp.MustCompile("[^a-z0-9]+")
	slug = reg.ReplaceAllString(slug, "-")

	// Remove leading and trailing hyphens
	slug = strings.Trim(slug, "-")

	// Replace multiple consecutive hyphens with a single one
	reg = regexp.MustCompile("-+")
	slug = reg.ReplaceAllString(slug, "-")

	// Limit length to 200 characters
	if len(slug) > 200 {
		slug = slug[:200]
		// Remove trailing hyphen if cut in the middle
		slug = strings.TrimRight(slug, "-")
	}

	return slug
}

// removeDiacritics removes accents and diacritics from text
func removeDiacritics(text string) string {
	// Romanian specific replacements
	replacements := map[rune]string{
		'ă': "a",
		'â': "a",
		'î': "i",
		'ș': "s",
		'ț': "t",
		'Ă': "a",
		'Â': "a",
		'Î': "i",
		'Ș': "s",
		'Ț': "t",
	}

	var result strings.Builder
	for _, char := range text {
		if replacement, found := replacements[char]; found {
			result.WriteString(replacement)
		} else {
			result.WriteRune(char)
		}
	}

	// Handle other Unicode diacritics
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	output, _, _ := transform.String(t, result.String())

	return output
}

// GenerateUniqueSlug appends a counter to make slug unique if needed
func GenerateUniqueSlug(baseSlug string, counter int) string {
	if counter == 0 {
		return baseSlug
	}
	return baseSlug + "-" + string(rune('0'+counter))
}
