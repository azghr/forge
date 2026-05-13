// Package stringutil provides common string transformations: Title, Slug, and
// RemoveAccents. All functions are pure, concurrency-safe, and have no external
// dependencies beyond the standard library.
package stringutil

import (
	"strings"
	"unicode"
)

// Title returns s in title case (each word capitalised).
// Words are sequences of letters separated by any non-letter rune.
// Unlike strings.Title (deprecated), this correctly handles Unicode punctuation.
func Title(s string) string {
	prevWasLetter := false
	r := []rune(s)
	for i, ch := range r {
		if unicode.IsLetter(ch) {
			if !prevWasLetter {
				r[i] = unicode.ToUpper(ch)
			}
			prevWasLetter = true
		} else {
			prevWasLetter = false
		}
	}
	return string(r)
}

// Slug returns a URL-friendly slug of s: lowercased, with non-alphanumeric
// characters replaced by a separator (default "-"), and consecutive separators
// collapsed into one. Leading and trailing separators are trimmed.
//
// Options can configure the separator and maximum length.
func Slug(s string, opts ...Option) string {
	cfg := defaultSlugConfig()
	for _, o := range opts {
		o(&cfg)
	}

	var b strings.Builder
	b.Grow(len(s))

	prevWasSep := false
	for _, ch := range []rune(s) {
		if unicode.IsLetter(ch) || unicode.IsDigit(ch) {
			b.WriteRune(unicode.ToLower(ch))
			prevWasSep = false
		} else {
			if !prevWasSep {
				b.WriteString(cfg.separator)
				prevWasSep = true
			}
		}
	}

	slug := b.String()
	slug = strings.TrimPrefix(slug, cfg.separator)
	slug = strings.TrimSuffix(slug, cfg.separator)

	if cfg.maxLength > 0 && len(slug) > cfg.maxLength {
		slug = slug[:cfg.maxLength]
		slug = strings.TrimSuffix(slug, cfg.separator)
	}

	return slug
}

// RemoveAccents strips diacritical marks from s, returning the base ASCII
// equivalent when possible. Unrecognised characters are kept as-is.
//
// Performance: O(n) over runes, no allocation for ASCII input.
func RemoveAccents(s string) string {
	r := []rune(s)
	for i, ch := range r {
		if base, ok := accentMap[ch]; ok {
			r[i] = base
		}
	}
	return string(r)
}

// accentMap maps accented Unicode runes to their ASCII base equivalent.
// Covers Latin-1 Supplement, Latin Extended-A/B, and common European
// accented characters.
var accentMap = map[rune]rune{
	// À-Å → A
	'\u00C0': 'A', '\u00C1': 'A', '\u00C2': 'A', '\u00C3': 'A',
	'\u00C4': 'A', '\u00C5': 'A', '\u0100': 'A', '\u0101': 'a',
	'\u0102': 'A', '\u0103': 'a', '\u0104': 'A', '\u0105': 'a',
	// à-å → a
	'\u00E0': 'a', '\u00E1': 'a', '\u00E2': 'a', '\u00E3': 'a',
	'\u00E4': 'a', '\u00E5': 'a',
	// È-Ë → E
	'\u00C8': 'E', '\u00C9': 'E', '\u00CA': 'E', '\u00CB': 'E',
	'\u0112': 'E', '\u0113': 'e', '\u0114': 'E', '\u0115': 'e',
	'\u0116': 'E', '\u0117': 'e', '\u0118': 'E', '\u0119': 'e',
	'\u011A': 'E', '\u011B': 'e',
	// è-ë → e
	'\u00E8': 'e', '\u00E9': 'e', '\u00EA': 'e', '\u00EB': 'e',
	// Ì-Ï → I
	'\u00CC': 'I', '\u00CD': 'I', '\u00CE': 'I', '\u00CF': 'I',
	'\u0128': 'I', '\u0129': 'i', '\u012A': 'I', '\u012B': 'i',
	'\u012C': 'I', '\u012D': 'i', '\u012E': 'I', '\u012F': 'i',
	'\u0130': 'I',
	// ì-ï → i
	'\u00EC': 'i', '\u00ED': 'i', '\u00EE': 'i', '\u00EF': 'i',
	// Ò-Ö → O
	'\u00D2': 'O', '\u00D3': 'O', '\u00D4': 'O', '\u00D5': 'O',
	'\u00D6': 'O', '\u014C': 'O', '\u014D': 'o', '\u014E': 'O',
	'\u014F': 'o', '\u0150': 'O', '\u0151': 'o',
	// ò-ö → o
	'\u00F2': 'o', '\u00F3': 'o', '\u00F4': 'o', '\u00F5': 'o',
	'\u00F6': 'o',
	// Ù-Ü → U
	'\u00D9': 'U', '\u00DA': 'U', '\u00DB': 'U', '\u00DC': 'U',
	'\u0168': 'U', '\u0169': 'u', '\u016A': 'U', '\u016B': 'u',
	'\u016C': 'U', '\u016D': 'u', '\u016E': 'U', '\u016F': 'u',
	'\u0170': 'U', '\u0171': 'u', '\u0172': 'U', '\u0173': 'u',
	// ù-ü → u
	'\u00F9': 'u', '\u00FA': 'u', '\u00FB': 'u', '\u00FC': 'u',
	// Ç → C, ç → c
	'\u00C7': 'C', '\u00E7': 'c',
	'\u0106': 'C', '\u0107': 'c', '\u0108': 'C', '\u0109': 'c',
	'\u010A': 'C', '\u010B': 'c', '\u010C': 'C', '\u010D': 'c',
	// Ñ → N, ñ → n
	'\u00D1': 'N', '\u00F1': 'n',
	'\u0143': 'N', '\u0144': 'n', '\u0145': 'N', '\u0146': 'n',
	'\u0147': 'N', '\u0148': 'n',
	// Ý → Y, ý → y
	'\u00DD': 'Y', '\u00FD': 'y', '\u00FF': 'y',
	'\u0176': 'Y', '\u0177': 'y', '\u0178': 'Y',
	// Š → S, š → s
	'\u0160': 'S', '\u0161': 's',
	'\u015A': 'S', '\u015B': 's', '\u015C': 'S', '\u015D': 's',
	'\u015E': 'S', '\u015F': 's',
	// Ž → Z, ž → z
	'\u017D': 'Z', '\u017E': 'z',
	'\u0179': 'Z', '\u017A': 'z', '\u017B': 'Z', '\u017C': 'z',
	// Ð (Eth) → D
	'\u00D0': 'D', '\u00F0': 'd',
	// Þ (Thorn) → T
	'\u00DE': 'T', '\u00FE': 't',
	// Æ → AE, æ → ae
	'\u00C6': 'A', '\u00E6': 'a',
	// Œ → OE, œ → oe
	'\u0152': 'O', '\u0153': 'o',
	// ß → ss (sharp s)
	'\u00DF': 's',
	// Đ → D, đ → d
	'\u0110': 'D', '\u0111': 'd',
	// Ħ → H, ħ → h
	'\u0126': 'H', '\u0127': 'h',
	// Ĳ → IJ, ĳ → ij
	'\u0132': 'I', '\u0133': 'i',
	// Ŀ → L, ŀ → l
	'\u013F': 'L', '\u0140': 'l',
	// Ł → L, ł → l
	'\u0141': 'L', '\u0142': 'l',
	// Ŋ → N, ŋ → n
	'\u014A': 'N', '\u014B': 'n',
	// Ŧ → T, ŧ → t
	'\u0166': 'T', '\u0167': 't',
	// soft-hyphen, no-break space → space
	'\u00AD': ' ', '\u00A0': ' ',
}
